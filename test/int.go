package tester

/*
Constants Generator
Help to auto-generated correct constants for lists
Martin Prestone (c) 2024
github.com/monopolly

Struct options
#index         Create int index anyway
#up            Uppercase for golang vars
#simple        Const only
name=Status    Named var. Ex name=Category
type=string    Type of const list. Ex type=string, int by default
#swift         Create swift enum

Fields options
#              Add lists for fields. Ex: id int //#readonly #must...
list{}         Create a.List() for field. Ex: list{colors, text}
iota=100       Add iota counter
go{}           Add new golang var name for fields. Ex go{ArtList}
title{}        Add title for fields. Ex title{News}
desc{}         Add desc for fields. Ex desc{New movie is good}
@              Custom values. Ex: id int //@code=401 @some=ErrSome
name{}         Add new name for fields. Ex name{newjsonname}
*/
type Status int

const (
	New     = Status(iota) + 1 //desc{Review} url{https:mice.com} @code=400 iota=1
	Mistake                    //desc{User must fix errors.} code{500} CodeName{Mistake}
	Review                     //desc{Waiting for the moderator review after user fix errors} #skip
	Active                     //desc{Active asset. User can get estimate changes}
	Public                     //desc{Asset is on the marketplace. Anyone could buy it}
	Ban                        //desc{Asset banned by moderator}
	Deleted                    //desc{Asset was deleted from platform. Nobody see it except moderators}
	Sold                       //desc{Asset hidden from platform. Nobody see it except moderators}
)

const (
	Transfered = Status(iota) + 101 //#virtual desc{Virtual status. Asset transfered to other user} iota=101
	Updated                         //#virtual desc{Virtual status. Asset has been updated}
	Locked                          //#virtual desc{Virtual status. Asset locked by moderator}
	Minted                          //#virtual desc{Virtual status. Mint tokens for asset}
	Certified                       //#virtual desc{Virtual status. Got Stima}
	Storage                         //#virtual desc{Virtual status. Move to storage}
	Insurance                       //#virtual desc{Virtual status. Create an insurance}
)

func (a Status) String() string {
	switch a {
	case New:
		return "new"
	case Mistake:
		return "mistake"
	case Review:
		return "review"
	case Active:
		return "active"
	case Public:
		return "public"
	case Ban:
		return "ban"
	case Deleted:
		return "deleted"
	case Sold:
		return "sold"
	case Transfered:
		return "transfered"
	case Updated:
		return "updated"
	case Locked:
		return "locked"
	case Minted:
		return "minted"
	case Certified:
		return "certified"
	case Storage:
		return "storage"
	case Insurance:
		return "insurance"
	default:
		return ""
	}
}

func (a Status) Title() string {
	switch a {
	case New:
		return "New"
	case Mistake:
		return "Mistake"
	case Review:
		return "Review"
	case Active:
		return "Active"
	case Public:
		return "Public"
	case Ban:
		return "Ban"
	case Deleted:
		return "Deleted"
	case Sold:
		return "Sold"
	case Transfered:
		return "Transfered"
	case Updated:
		return "Updated"
	case Locked:
		return "Locked"
	case Minted:
		return "Minted"
	case Certified:
		return "Certified"
	case Storage:
		return "Storage"
	case Insurance:
		return "Insurance"
	default:
		return ""
	}
}

func (a Status) Desc() string {
	switch a {
	case New:
		return "Review"
	case Mistake:
		return "User must fix errors."
	case Review:
		return "Waiting for the moderator review after user fix errors"
	case Active:
		return "Active asset. User can get estimate changes"
	case Public:
		return "Asset is on the marketplace. Anyone could buy it"
	case Ban:
		return "Asset banned by moderator"
	case Deleted:
		return "Asset was deleted from platform. Nobody see it except moderators"
	case Sold:
		return "Asset hidden from platform. Nobody see it except moderators"
	case Transfered:
		return "Virtual status. Asset transfered to other user"
	case Updated:
		return "Virtual status. Asset has been updated"
	case Locked:
		return "Virtual status. Asset locked by moderator"
	case Minted:
		return "Virtual status. Mint tokens for asset"
	case Certified:
		return "Virtual status. Got Stima"
	case Storage:
		return "Virtual status. Move to storage"
	case Insurance:
		return "Virtual status. Create an insurance"
	default:
		return ""
	}
}

func (a Status) Int() int {
	return int(a)
}

func ValidIndexType(v Status) bool {
	switch v {
	case New, Mistake, Review, Active, Public, Ban, Deleted, Sold, Transfered, Updated, Locked, Minted, Certified, Storage, Insurance:
		return true
	default:
		return false
	}
}

func (a Status) Valid() bool {
	return ValidIndexType(a)
}

func Valid(v int) bool {
	switch Status(v) {
	case New, Mistake, Review, Active, Public, Ban, Deleted, Sold, Transfered, Updated, Locked, Minted, Certified, Storage, Insurance:
		return true
	default:
		return false
	}
}

func StatusIndexes() []Status {
	return []Status{New, Mistake, Review, Active, Public, Ban, Deleted, Sold, Transfered, Updated, Locked, Minted, Certified, Storage, Insurance}
}

func StatusSkipIndexes() []Status {
	return []Status{Review}
}

func (a Status) IsSkip() bool {
	switch a {
	case Review:
		return true
	default:
		return false
	}
}
func StatusVirtualIndexes() []Status {
	return []Status{Transfered, Updated, Locked, Minted, Certified, Storage, Insurance}
}

func (a Status) IsVirtual() bool {
	switch a {
	case Transfered, Updated, Locked, Minted, Certified, Storage, Insurance:
		return true
	default:
		return false
	}
}

func Index(v string) int {
	switch v {
	case "new":
		return 1
	case "mistake":
		return 2
	case "review":
		return 3
	case "active":
		return 4
	case "public":
		return 5
	case "ban":
		return 6
	case "deleted":
		return 7
	case "sold":
		return 8
	case "transfered":
		return 101
	case "updated":
		return 102
	case "locked":
		return 103
	case "minted":
		return 104
	case "certified":
		return 105
	case "storage":
		return 106
	case "insurance":
		return 107
	default:
		return 0
	}
}

func StatusSkipList() []Status {
	return []Status{Review}
}

func (a Status) Skip() bool {
	switch a {
	case Review:
		return true
	default:
		return false
	}
}
func StatusVirtualList() []Status {
	return []Status{Transfered, Updated, Locked, Minted, Certified, Storage, Insurance}
}

func (a Status) Virtual() bool {
	switch a {
	case Transfered, Updated, Locked, Minted, Certified, Storage, Insurance:
		return true
	default:
		return false
	}
}

func (a Status) url() string {
	switch a {
	case New:
		return "https:mice.com"
	default:
		return ""
	}
}

func (a Status) code() int {
	switch a {
	case Mistake:
		return 500
	default:
		return 0
	}
}

func (a Status) CodeName() string {
	switch a {
	case Mistake:
		return "Mistake"
	default:
		return ""
	}
}
