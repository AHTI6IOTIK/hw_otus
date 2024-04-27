package hw09structvalidator

import (
	"errors"
	"fmt"
	"regexp"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

var (
	ErrLen        = errors.New("the string length error")
	ErrRegexp     = errors.New("the string does not match the regular expression")
	ErrList       = errors.New("the string is not included in the list")
	ErrNumber     = errors.New("number validation error")
	ErrNumberList = errors.New("the number is not included in the list")
)

type Constraint[T comparable] interface {
	Error() error
	isValid(value T) bool
}

type ConstraintWrapper[T comparable] struct {
	constraints []Constraint[T]
}

func (c *ConstraintWrapper[T]) add(constraint Constraint[T]) {
	c.constraints = append(c.constraints, constraint)
}

//nolint:unused
func (c *ConstraintWrapper[T]) isValid(value T) bool {
	for _, c2 := range c.constraints {
		if !c2.isValid(value) {
			return false
		}
	}

	return true
}

func (c *ConstraintWrapper[T]) Error() error {
	var result error

	for _, c2 := range c.constraints {
		if result == nil {
			result = c2.Error()
		} else {
			result = fmt.Errorf("%w | %v", result, c2.Error()) //nolint:errorlint,nolintlint
		}
	}

	return result
}

type LenConstraint[T ~string | ~[]string] struct {
	constraint int
}

func NewLenConstraint[T ~string | ~[]string](v int) *LenConstraint[T] {
	return &LenConstraint[T]{
		constraint: v,
	}
}

//nolint:unused
func (l *LenConstraint[T]) isValid(value T) bool {
	return len(value) == l.constraint
}

func (l *LenConstraint[T]) Error() error {
	return fmt.Errorf("%w | should be: %d", ErrLen, l.constraint)
}

type RegexpConstraint[T ~string] struct {
	constraint *regexp.Regexp
}

func NewRegexpConstraint[T ~string](v *regexp.Regexp) *RegexpConstraint[T] {
	return &RegexpConstraint[T]{
		constraint: v,
	}
}

//nolint:unused
func (r *RegexpConstraint[T]) isValid(value T) bool {
	return r.constraint.Match([]byte(value))
}

func (r *RegexpConstraint[T]) Error() error {
	return fmt.Errorf("%w | %s", ErrRegexp, r.constraint.String())
}

type InStringsConstraint[T string] struct {
	constraint []T
}

func NewInStringsConstraint[T string](v []T) *InStringsConstraint[T] {
	return &InStringsConstraint[T]{
		constraint: v,
	}
}

//nolint:unused
func (i *InStringsConstraint[T]) isValid(value T) bool {
	return slices.Contains[[]T, T](i.constraint, value)
}

func (i *InStringsConstraint[T]) Error() error {
	return fmt.Errorf("%w | %v", ErrList, i.constraint)
}

type MinConstraint[T constraints.Integer] struct {
	constraint T
}

func NewMinConstraint[T constraints.Integer](v T) *MinConstraint[T] {
	return &MinConstraint[T]{
		constraint: v,
	}
}

//nolint:unused
func (m *MinConstraint[T]) isValid(value T) bool {
	return value > m.constraint
}

func (m *MinConstraint[T]) Error() error {
	return fmt.Errorf("%w  | the number must be greater: %d", ErrNumber, m.constraint)
}

type MaxConstraint[T constraints.Integer] struct {
	constraint T
}

func NewMaxConstraint[T constraints.Integer](v T) *MaxConstraint[T] {
	return &MaxConstraint[T]{
		constraint: v,
	}
}

//nolint:unused
func (m *MaxConstraint[T]) isValid(value T) bool {
	return value <= m.constraint
}

func (m *MaxConstraint[T]) Error() error {
	return fmt.Errorf("%w  | the number must be less than: %d", ErrNumber, m.constraint)
}

type InIntConstraint[T constraints.Integer] struct {
	constraint []T
}

func NewInIntConstraint[T constraints.Integer](v []T) *InIntConstraint[T] {
	return &InIntConstraint[T]{
		constraint: v,
	}
}

//nolint:unused
func (i *InIntConstraint[T]) isValid(value T) bool {
	return slices.Contains(i.constraint, value)
}

func (i *InIntConstraint[T]) Error() error {
	return fmt.Errorf("%w  | the number is not included in the list: %v", ErrNumberList, i.constraint)
}
