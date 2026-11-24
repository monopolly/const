//enum list
enum Status: Int {

	//case
	case new = 1, mistake, review, active, public, ban, deleted, sold, transfered = 101, updated, locked, minted, certified, storage, insurance

	//id
	var id: Int {
		get {
			switch self {
				case .new:         return 1
				case .mistake:     return 2
				case .review:      return 3
				case .active:      return 4
				case .public:      return 5
				case .ban:         return 6
				case .deleted:     return 7
				case .sold:        return 8
				case .transfered:  return 101
				case .updated:     return 102
				case .locked:      return 103
				case .minted:      return 104
				case .certified:   return 105
				case .storage:     return 106
				case .insurance:   return 107
			}
		}
	}
	
	//title
	var title: String {
		get {
			switch self {
				case .new:         return "New"
				case .mistake:     return "Mistake"
				case .review:      return "Review"
				case .active:      return "Active"
				case .public:      return "Public"
				case .ban:         return "Ban"
				case .deleted:     return "Deleted"
				case .sold:        return "Sold"
				case .transfered:  return "Transfered"
				case .updated:     return "Updated"
				case .locked:      return "Locked"
				case .minted:      return "Minted"
				case .certified:   return "Certified"
				case .storage:     return "Storage"
				case .insurance:   return "Insurance"
			}
		}
	}
	
	//desc
	var desc: String {
		get {
			switch self {
				case .new:         return "Review"
				case .mistake:     return "User must fix errors."
				case .review:      return "Waiting for the moderator review after user fix errors"
				case .active:      return "Active asset. User can get estimate changes"
				case .public:      return "Asset is on the marketplace. Anyone could buy it"
				case .ban:         return "Asset banned by moderator"
				case .deleted:     return "Asset was deleted from platform. Nobody see it except moderators"
				case .sold:        return "Asset hidden from platform. Nobody see it except moderators"
				case .transfered:  return "Virtual status. Asset transfered to other user"
				case .updated:     return "Virtual status. Asset has been updated"
				case .locked:      return "Virtual status. Asset locked by moderator"
				case .minted:      return "Virtual status. Mint tokens for asset"
				case .certified:   return "Virtual status. Got Stima"
				case .storage:     return "Virtual status. Move to storage"
				case .insurance:   return "Virtual status. Create an insurance"
			}
		}
	}
	
	var skip: Bool {
		get {
			switch self {
				case .review: return true
				default: return false
			}
		}
	}
	
	var virtual: Bool {
		get {
			switch self {
				case .transfered: return true
				case .updated: return true
				case .locked: return true
				case .minted: return true
				case .certified: return true
				case .storage: return true
				case .insurance: return true
				default: return false
			}
		}
	}
	
}
