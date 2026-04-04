package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ─── Company ──────────────────────────────────────────────────────────────────

func (Company) TableName() string { return "companies" }

type Company struct {
	Id        uuid.UUID       `gorm:"type:uuid;primaryKey;default:gen_random_uuid();" json:"id"`
	Name      string          `gorm:"type:character varying(150);not null;" json:"name"`
	LegalName string          `gorm:"type:character varying(200);" json:"legal_name"`
	TaxID     string          `gorm:"type:character varying(50);" json:"tax_id"`
	Address   string          `gorm:"type:text;" json:"address"`
	Phone     string          `gorm:"type:character varying(30);" json:"phone"`
	Email     string          `gorm:"type:character varying(150);" json:"email"`
	Industry  string          `gorm:"type:character varying(100);" json:"industry"`
	Currency  string          `gorm:"type:character varying(10);not null;default:'IDR';" json:"currency"`
	Status    int             `gorm:"type:int;not null;default:1;" json:"status"` // 1=active, 0=inactive
	DeletedAt *time.Time      `gorm:"default:null;" json:"deleted_at,omitempty"`
	CreatedBy uuid.UUID       `gorm:"type:uuid;not null;" json:"created_by"`
	Members   []CompanyMember `gorm:"foreignKey:CompanyId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"members,omitempty"`
	CreatedAt time.Time       `gorm:"autoCreateTime;" json:"created_at"`
	UpdatedAt time.Time       `gorm:"autoUpdateTime;" json:"updated_at"`
}

func (Company) SearchableFields() []string {
	return []string{"name", "legal_name", "email", "tax_id"}
}

func (Company) ApplyFilters(db *gorm.DB, filters map[string][]string) *gorm.DB {
	if values, ok := filters["status"]; ok {
		db = db.Where("companies.status IN ?", values)
	}
	return db
}

// ─── CompanyRole ──────────────────────────────────────────────────────────────

// CompanyRole defines what a user can do within a company.
// Hierarchy (descending): superadmin > owner > admin > accountant > viewer
// superadmin is platform-level and set manually in DB — never via API.
type CompanyRole string

const (
	RoleSuperadmin CompanyRole = "superadmin" // platform-wide, manual DB only
	RoleOwner      CompanyRole = "owner"      // full company access, can delete company
	RoleAdmin      CompanyRole = "admin"      // manage users, settings, fiscal, COA
	RoleAccountant CompanyRole = "accountant" // create/edit/post/void transactions
	RoleViewer     CompanyRole = "viewer"     // read-only
)

// PermissionMap defines what each role can do.
// Used by the RBAC middleware to gate endpoints.
var PermissionMap = map[CompanyRole][]string{
	RoleSuperadmin: {"*"}, // all permissions
	RoleOwner: {
		"company:read", "company:update", "company:delete",
		"users:read", "users:manage",
		"coa:read", "coa:manage",
		"fiscal:read", "fiscal:manage",
		"transactions:read", "transactions:write", "transactions:post",
		"reports:read",
	},
	RoleAdmin: {
		"company:read", "company:update",
		"users:read", "users:manage",
		"coa:read", "coa:manage",
		"fiscal:read", "fiscal:manage",
		"transactions:read", "transactions:write", "transactions:post",
		"reports:read",
	},
	RoleAccountant: {
		"company:read",
		"users:read",
		"coa:read",
		"fiscal:read",
		"transactions:read", "transactions:write", "transactions:post",
		"reports:read",
	},
	RoleViewer: {
		"company:read",
		"coa:read",
		"fiscal:read",
		"transactions:read",
		"reports:read",
	},
}

// HasPermission checks if a role has a specific permission.
func HasPermission(role CompanyRole, permission string) bool {
	perms, ok := PermissionMap[role]
	if !ok {
		return false
	}
	for _, p := range perms {
		if p == "*" || p == permission {
			return true
		}
	}
	return false
}

// ─── CompanyMember ────────────────────────────────────────────────────────────

func (CompanyMember) TableName() string { return "company_members" }

// CompanyMember is the join table between User and Company.
// A user can belong to multiple companies with different roles.
type CompanyMember struct {
	Id        uuid.UUID   `gorm:"type:uuid;primaryKey;default:gen_random_uuid();" json:"id"`
	CompanyId uuid.UUID   `gorm:"type:uuid;not null;index;" json:"company_id"`
	UserId    uuid.UUID   `gorm:"type:uuid;not null;index;" json:"user_id"`
	Role      CompanyRole `gorm:"type:character varying(20);not null;" json:"role"`
	InvitedBy uuid.UUID   `gorm:"type:uuid;" json:"invited_by"`
	JoinedAt  time.Time   `gorm:"autoCreateTime;" json:"joined_at"`
	UpdatedAt time.Time   `gorm:"autoUpdateTime;" json:"updated_at"`

	// Associations (preload only, not stored)
	Company Company `gorm:"foreignKey:CompanyId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"company,omitempty"`
	User    User    `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
}
