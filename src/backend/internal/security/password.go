package security

import (
	"errors"
	"regexp"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

var (
	// Cost bcrypt cost 值
	Cost = 10

	// ErrPasswordTooShort 密码太短
	ErrPasswordTooShort = errors.New("password length must be at least 8 characters")
	// ErrPasswordNoUppercase 缺少大写字母
	ErrPasswordNoUppercase = errors.New("password must contain at least one uppercase letter")
	// ErrPasswordNoLowercase 缺少小写字母
	ErrPasswordNoLowercase = errors.New("password must contain at least one lowercase letter")
	// ErrPasswordNoDigit 缺少数字
	ErrPasswordNoDigit = errors.New("password must contain at least one digit")
	// ErrPasswordNoSpecialChar 缺少特殊字符
	ErrPasswordNoSpecialChar = errors.New("password must contain at least one special character")
)

// PasswordComplexityError 密码复杂度错误
type PasswordComplexityError struct {
	MissingRequirements []error
}

// Error 实现 error 接口
func (e *PasswordComplexityError) Error() string {
	if len(e.MissingRequirements) == 0 {
		return "password complexity validation failed"
	}

	msg := "password does not meet complexity requirements: "
	for i, err := range e.MissingRequirements {
		if i > 0 {
			msg += "; "
		}
		msg += err.Error()
	}
	return msg
}

// HashPassword 对密码进行哈希
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), Cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword 验证密码
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ValidatePasswordComplexity 校验密码复杂度
func ValidatePasswordComplexity(password string) error {
	var missingRequirements []error

	// 校验长度至少8位
	if len(password) < 8 {
		missingRequirements = append(missingRequirements, ErrPasswordTooShort)
	}

	// 校验包含大写字母
	hasUppercase := false
	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUppercase = true
			break
		}
	}
	if !hasUppercase {
		missingRequirements = append(missingRequirements, ErrPasswordNoUppercase)
	}

	// 校验包含小写字母
	hasLowercase := false
	for _, char := range password {
		if unicode.IsLower(char) {
			hasLowercase = true
			break
		}
	}
	if !hasLowercase {
		missingRequirements = append(missingRequirements, ErrPasswordNoLowercase)
	}

	// 校验包含数字
	hasDigit := false
	for _, char := range password {
		if unicode.IsDigit(char) {
			hasDigit = true
			break
		}
	}
	if !hasDigit {
		missingRequirements = append(missingRequirements, ErrPasswordNoDigit)
	}

	// 校验包含特殊字符
	specialCharPattern := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?~]`)
	if !specialCharPattern.MatchString(password) {
		missingRequirements = append(missingRequirements, ErrPasswordNoSpecialChar)
	}

	// 如果有缺失的要求，返回错误
	if len(missingRequirements) > 0 {
		return &PasswordComplexityError{
			MissingRequirements: missingRequirements,
		}
	}

	return nil
}
