package guard

type Permissions []int

const (
	DeviceAccept = iota + 1
	DeviceReject
	DeviceRemove
	DeviceConnect
	DeviceRename
	DeviceDetails

	DeviceCreateTag
	DeviceUpdateTag
	DeviceRemoveTag
	DeviceRenameTag
	DeviceDeleteTag

	SessionPlay
	SessionClose
	SessionRemove
	SessionDetails

	FirewallCreate
	FirewallEdit
	FirewallRemove

	FirewallAddTag
	FirewallRemoveTag
	FirewallUpdateTag

	PublicKeyCreate
	PublicKeyEdit
	PublicKeyRemove

	PublicKeyAddTag
	PublicKeyRemoveTag
	PublicKeyUpdateTag

	NamespaceRename
	NamespaceAddMember
	NamespaceRemoveMember
	NamespaceEditMember
	NamespaceEnableSessionRecord
	NamespaceDelete

	BillingChooseDevices
	BillingAddPaymentMethod
	BillingUpdatePaymentMethod
	BillingRemovePaymentMethod
	BillingCancelSubscription
	BillingCreateSubscription
	BillingGetPaymentMethod
	BillingGetSubscription
)

var observerPermissions = Permissions{
	DeviceConnect,
	DeviceDetails,
	SessionDetails,
}

var operatorPermissions = Permissions{
	DeviceAccept,
	DeviceReject,
	DeviceConnect,
	DeviceRename,
	DeviceDetails,

	DeviceCreateTag,
	DeviceUpdateTag,
	DeviceRemoveTag,
	DeviceRenameTag,
	DeviceDeleteTag,

	SessionDetails,
}

var adminPermissions = Permissions{
	DeviceAccept,
	DeviceReject,
	DeviceRemove,
	DeviceConnect,
	DeviceRename,
	DeviceDetails,

	DeviceCreateTag,
	DeviceUpdateTag,
	DeviceRemoveTag,
	DeviceRenameTag,
	DeviceDeleteTag,

	SessionPlay,
	SessionClose,
	SessionRemove,
	SessionDetails,

	FirewallCreate,
	FirewallEdit,
	FirewallRemove,
	FirewallAddTag,
	FirewallRemoveTag,
	FirewallUpdateTag,

	PublicKeyCreate,
	PublicKeyEdit,
	PublicKeyRemove,
	PublicKeyAddTag,
	PublicKeyRemoveTag,
	PublicKeyUpdateTag,

	NamespaceRename,
	NamespaceAddMember,
	NamespaceRemoveMember,
	NamespaceEditMember,
	NamespaceEnableSessionRecord,
}

var ownerPermissions = Permissions{
	DeviceAccept,
	DeviceReject,
	DeviceRemove,
	DeviceConnect,
	DeviceRename,
	DeviceDetails,

	DeviceCreateTag,
	DeviceUpdateTag,
	DeviceRemoveTag,
	DeviceRenameTag,
	DeviceDeleteTag,

	SessionPlay,
	SessionClose,
	SessionRemove,
	SessionDetails,

	FirewallCreate,
	FirewallEdit,
	FirewallRemove,
	FirewallAddTag,
	FirewallRemoveTag,
	FirewallUpdateTag,

	PublicKeyCreate,
	PublicKeyEdit,
	PublicKeyRemove,
	PublicKeyAddTag,
	PublicKeyRemoveTag,
	PublicKeyUpdateTag,

	NamespaceRename,
	NamespaceAddMember,
	NamespaceRemoveMember,
	NamespaceEditMember,
	NamespaceEnableSessionRecord,
	NamespaceDelete,

	BillingChooseDevices,
	BillingAddPaymentMethod,
	BillingUpdatePaymentMethod,
	BillingRemovePaymentMethod,
	BillingCancelSubscription,
	BillingCreateSubscription,
	BillingGetSubscription,
}
