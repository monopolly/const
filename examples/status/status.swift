//enum list
enum Status: Int {

	//case
	case new = 1, mistake = 2, review = 3, active = 4, `public` = 5, ban = 6, deleted = 7, sold = 8, transferred = 101, updated = 102, locked = 103

	//id
	var id: Int {
		get {
			switch self {
				case .new:          return 1
				case .mistake:      return 2
				case .review:       return 3
				case .active:       return 4
				case .`public`:     return 5
				case .ban:          return 6
				case .deleted:      return 7
				case .sold:         return 8
				case .transferred:  return 101
				case .updated:      return 102
				case .locked:       return 103
			}
		}
	}
	
	//title
	var title: String {
		get {
			switch self {
				case .new:          return "New"
				case .mistake:      return "Mistake"
				case .review:       return "Review"
				case .active:       return "Active"
				case .`public`:     return "Public"
				case .ban:          return "Ban"
				case .deleted:      return "Deleted"
				case .sold:         return "Sold"
				case .transferred:  return "Transferred"
				case .updated:      return "Updated"
				case .locked:       return "Locked"
			}
		}
	}
	
	//desc
	var desc: String {
		get {
			switch self {
				case .new:          return "Waiting for review"
				case .mistake:      return "User must fix errors"
				case .review:       return "Waiting for the moderator review after user fix errors"
				case .active:       return "Active asset"
				case .`public`:     return "Asset is on the marketplace"
				case .ban:          return "Asset banned by moderator"
				case .deleted:      return "Asset was deleted"
				case .sold:         return "Asset was sold"
				case .transferred:  return "Asset transferred to other user"
				case .updated:      return "Asset has been updated"
				case .locked:       return "Asset locked by moderator"
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
				case .transferred: return true
				case .updated: return true
				case .locked: return true
				default: return false
			}
		}
	}
	
}
