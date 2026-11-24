//enum list
enum Category: String {

	//case
	case artwork = "artwork", bag = "bag", boat = "boat", car = "car", estate = "estate", jet = "jet", watch = "watch", wineSpirits = "wine & spirits"

	//id
	var id: String {
		get {
			switch self {
				case .artwork:      return "artwork"
				case .bag:          return "bag"
				case .boat:         return "boat"
				case .car:          return "car"
				case .estate:       return "estate"
				case .jet:          return "jet"
				case .watch:        return "watch"
				case .wineSpirits:  return "wine & spirits"
			}
		}
	}
	
	//title
	var title: String {
		get {
			switch self {
				case .artwork:      return "Art"
				case .bag:          return "Handbags"
				case .boat:         return "Yachts & Boats"
				case .car:          return "Cars"
				case .estate:       return "Real Estate"
				case .jet:          return "Airplanes"
				case .watch:        return "Watches"
				case .wineSpirits:  return "Wine & Spirits"
			}
		}
	}
	
	//desc
	var desc: String {
		get {
			switch self {
				case .estate:       return "Houses, flats and land"
				case .jet:          return "Business jets"
				default:            return ""
			}
		}
	}
	
	var active: Bool {
		get {
			switch self {
				case .artwork: return true
				case .bag: return true
				case .car: return true
				case .watch: return true
				default: return false
			}
		}
	}
	
}
