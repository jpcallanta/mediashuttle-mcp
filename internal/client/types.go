package client

import "time"

// Portal represents a Media Shuttle portal.
type Portal struct {
	ID             string          `json:"id,omitempty"`
	Name           string          `json:"name,omitempty"`
	URL            string          `json:"url,omitempty"`
	Type           string          `json:"type,omitempty"`
	CreatedOn      *time.Time      `json:"createdOn,omitempty"`
	LastModifiedOn *time.Time      `json:"lastModifiedOn,omitempty"`
	Authentication *Authentication `json:"authentication,omitempty"`
	LinkAuth       *LinkAuth       `json:"linkAuthentication,omitempty"`
	NotifyMembers  *bool           `json:"notifyMembers,omitempty"`
	ContactForm    *ContactForm    `json:"contactForm,omitempty"`
	SAMLHomeFolder string          `json:"samlHomeFolder,omitempty"`
}

// Authentication holds portal auth settings.
type Authentication struct {
	MediaShuttle bool `json:"mediaShuttle,omitempty"`
	SAML         bool `json:"saml,omitempty"`
}

// LinkAuth holds link-based auth settings.
type LinkAuth struct {
	AllowUnauthUploads   *bool `json:"allowUnauthenticatedUploads,omitempty"`
	AllowUnauthDownloads *bool `json:"allowUnauthenticatedDownloads,omitempty"`
}

// ContactForm holds custom contact form config.
type ContactForm struct {
	URL    string `json:"url,omitempty"`
	Width  *int   `json:"width,omitempty"`
	Height *int   `json:"height,omitempty"`
}

// UpdatePortal is the request body for PATCH /portals/{id}.
type UpdatePortal struct {
	Name           string          `json:"name,omitempty"`
	URL            string          `json:"url,omitempty"`
	Authentication *Authentication `json:"authentication,omitempty"`
	LinkAuth       *LinkAuth       `json:"linkAuthentication,omitempty"`
	NotifyMembers  *bool           `json:"notifyMembers,omitempty"`
	ContactForm    *ContactForm    `json:"contactForm,omitempty"`
	SAMLHomeFolder string          `json:"samlHomeFolder,omitempty"`
}

// PortalList is the response from listing portals.
type PortalList struct {
	Items []Portal `json:"items"`
}

// PortalMember represents a user on a portal.
type PortalMember struct {
	Email     string             `json:"email"`
	Role      string             `json:"role,omitempty"`
	ExpiresOn string             `json:"expiresOn,omitempty"`
	Perms     *PortalPermissions `json:"portalPermissions,omitempty"`
}

// PortalPermissions holds a member's portal permissions.
type PortalPermissions struct {
	CanSendToMember    *bool          `json:"canSendToMember,omitempty"`
	CanSendToNonMember *bool          `json:"canSendToNonMember,omitempty"`
	CanReceive         *bool          `json:"canReceive,omitempty"`
	CanSendFromShare   *bool          `json:"canSendFromShare,omitempty"`
	CanAutoDeliver     *bool          `json:"canDeliverAutomatically,omitempty"`
	CanSubmit          *bool          `json:"canSubmit,omitempty"`
	Folders            []PortalFolder `json:"folders,omitempty"`
}

// PortalFolder is a folder on a share portal.
type PortalFolder struct {
	Path string `json:"path,omitempty"`
}

// PortalMemberResponse is the response from listing portal members.
type PortalMemberResponse struct {
	Items []PortalMemberListItem `json:"items"`
	Error *APIError              `json:"error,omitempty"`
}

// PortalMemberListItem is a summary of a portal member.
type PortalMemberListItem struct {
	Email       string `json:"email"`
	LastLoginOn string `json:"lastLoginOn,omitempty"`
}

// ResponseForPortalMember is a full portal member response.
type ResponseForPortalMember struct {
	Email     string             `json:"email"`
	Role      string             `json:"role,omitempty"`
	ExpiresOn string             `json:"expiresOn,omitempty"`
	LastLogin string             `json:"lastLoginOn,omitempty"`
	Perms     *PortalPermissions `json:"portalPermissions,omitempty"`
}

// PortalStorage associates storage with a portal.
type PortalStorage struct {
	StorageID      string `json:"storageId"`
	RepositoryPath string `json:"repositoryPath,omitempty"`
}

// PortalStorageList is the response from listing portal storage.
type PortalStorageList struct {
	Items []PortalStorage `json:"items"`
}

// Storage is a storage location registered to an account.
type Storage struct {
	ID            string         `json:"id,omitempty"`
	Type          string         `json:"type,omitempty"`
	Relays        []string       `json:"relays,omitempty"`
	Status        string         `json:"status,omitempty"`
	Configuration map[string]any `json:"configuration,omitempty"`
}

// StorageList is the response from listing storage.
type StorageList struct {
	Items []Storage `json:"items"`
}

// Transfer is an active file transfer.
type Transfer struct {
	ID              string                `json:"id,omitempty"`
	PortalID        string                `json:"portalId,omitempty"`
	PackageID       string                `json:"packageId,omitempty"`
	State           string                `json:"state,omitempty"`
	Protocol        string                `json:"protocol,omitempty"`
	ConnectedServer string                `json:"connectedServer,omitempty"`
	Direction       string                `json:"direction,omitempty"`
	User            *TransferUser         `json:"user,omitempty"`
	StartTime       string                `json:"startTime,omitempty"`
	Details         *ActiveTransferDetail `json:"activeTransferDetails,omitempty"`
}

// TransferUser is the user performing a transfer.
type TransferUser struct {
	Email string `json:"email,omitempty"`
}

// ActiveTransferDetail holds live transfer stats.
type ActiveTransferDetail struct {
	EstTimeRemaining float64 `json:"estimatedTimeRemainingInSeconds,omitempty"`
	TransferRate     float64 `json:"transferRateInBitsPerSecond,omitempty"`
	CurrentFile      *File   `json:"currentFile,omitempty"`
}

// File is a file being transferred.
type File struct {
	Name string `json:"name,omitempty"`
	Size int64  `json:"size,omitempty"`
}

// TransferList is the response from listing transfers.
type TransferList struct {
	Items []Transfer `json:"items"`
}

// APIError is a standard API error response.
type APIError struct {
	StatusCode int    `json:"statusCode,omitempty"`
	Error      string `json:"error,omitempty"`
	Message    string `json:"message,omitempty"`
}
