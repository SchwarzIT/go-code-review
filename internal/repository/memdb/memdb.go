package memdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

const (
	default_db_filename = "coupons.data.json"
	baseDir             = "./data"
)

type Coupon struct {
	ID             string `json:"id"`
	Code           string `json:"code"`
	Discount       int    `json:"discount"`
	MinBasketValue int    `json:"min_basket_value"`
}

// Repository defines the in-memory storage for Coupons.
// It implements the repository interface.
type Repository struct {
	entries  map[string]*Coupon
	filePath string
	mu       sync.RWMutex
}

// NewRepository creates and returns a new Repository instance.
// filePath to store the db
func NewRepository(filePath string) (*Repository, error) {
	repo := &Repository{
		entries:  make(map[string]*Coupon),
		filePath: filePath,
	}

	// Load existing coupons from the file
	err := repo.loadFromFile()
	if err != nil {
		return nil, fmt.Errorf("failed to load data from file: %w", err)
	}

	return repo, nil
}

// NewRepositoryDefault creates and returns a new Repository instance.
// default file path value will be './data/coupons.data.json'
// otherwise, it initializes an empty repository and creates the file.
func NewRepositoryDefault() (*Repository, error) {
	repo := &Repository{
		entries:  make(map[string]*Coupon),
		filePath: filepath.Join(baseDir, default_db_filename),
	}

	// Load existing coupons from the file
	err := repo.loadFromFile()
	if err != nil {
		return nil, fmt.Errorf("failed to load data from file: %w", err)
	}

	return repo, nil
}

// RepositoryInterface defines the methods that the Repository implements.
// Exported for external usage if needed.
type RepositoryInterface interface {
	FindByCode(string) (*Coupon, error)
	Save(*Coupon) error
}

// Custom errors for better error handling.
var (
	ErrCouponNotFound = errors.New("coupon not found")
	ErrInvalidCoupon  = errors.New("invalid coupon")
)

// FindByCode retrieves a Coupon by its code.
// It returns a copy of the Coupon to prevent external modifications.
func (r *Repository) FindByCode(code string) (Coupon, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	coupon, exists := r.entries[code]
	if !exists {
		return Coupon{}, ErrCouponNotFound
	}

	// Return a copy to maintain immutability.
	return *coupon, nil
}

// loadFromFile reads the coupons from the data.json file into the repository.
// If the file does not exist, it initializes an empty repository and creates the file.
func (r *Repository) loadFromFile() (err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	log.Printf("Loading data from '%s'", r.filePath)
	file, err := os.Open(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// File does not exist; create an empty new one
			file, createErr := os.Create(r.filePath)
			if createErr != nil {
				return fmt.Errorf("unable to create data file: %w", createErr)
			}
			defer func() {
				if cerr := file.Close(); cerr != nil && err != nil {
					err = cerr
				}
			}()
			return nil
		}
		return fmt.Errorf("unable to open data file: %w", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil && err != nil {
			err = cerr
		}
	}()

	decoder := json.NewDecoder(file)
	var coupons []*Coupon
	if err := decoder.Decode(&coupons); err != nil {
		if err.Error() == "EOF" {
			return nil
		}
		return fmt.Errorf("error decoding JSON from file: %w", err)
	}

	for _, coupon := range coupons {
		if coupon != nil && coupon.Code != "" {
			couponCopy := *coupon
			r.entries[coupon.Code] = &couponCopy
		}
	}

	return nil
}

// Save stores a Coupon in the repository.
// It returns an error if the coupon is nil or has an empty code.
func (r *Repository) Save(coupon *Coupon) error {
	if coupon == nil {
		return ErrInvalidCoupon
	}
	if coupon.Code == "" {
		return fmt.Errorf("%w: coupon code is empty", ErrInvalidCoupon)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Store a copy to prevent external modifications affecting the repository.
	couponCopy := *coupon
	r.entries[coupon.Code] = &couponCopy

	// Persist the updated entries to the file
	if err := r.saveToFile(); err != nil {
		return fmt.Errorf("failed to save coupon to file: %w", err)
	}

	return nil
}

// saveToFile writes the current state of coupons to the data.json file.
// It directly writes to data.json without using a temporary file.
func (r *Repository) saveToFile() (err error) {
	// Prepare a slice to hold coupons for JSON encoding
	coupons := make([]*Coupon, 0, len(r.entries))
	for _, coupon := range r.entries {
		coupons = append(coupons, coupon)
	}

	// Open the file with write permissions, create it if it doesn't exist, truncate it
	file, err := os.OpenFile(r.filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("unable to open data file for writing: %w", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil && err != nil {
			err = cerr
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // For pretty-printing

	if err := encoder.Encode(coupons); err != nil {
		return fmt.Errorf("error encoding JSON to file: %w", err)
	}

	return nil
}
