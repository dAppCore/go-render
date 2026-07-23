package php

func TestPackages_LinkPackages_Good(t *T) {
	subject := LinkPackages
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestPackages_LinkPackages_Bad(t *T) {
	subject := LinkPackages
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestPackages_LinkPackages_Ugly(t *T) {
	subject := LinkPackages
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestPackages_UnlinkPackages_Good(t *T) {
	subject := UnlinkPackages
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestPackages_UnlinkPackages_Bad(t *T) {
	subject := UnlinkPackages
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestPackages_UnlinkPackages_Ugly(t *T) {
	subject := UnlinkPackages
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestPackages_UpdatePackages_Good(t *T) {
	subject := UpdatePackages
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestPackages_UpdatePackages_Bad(t *T) {
	subject := UpdatePackages
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestPackages_UpdatePackages_Ugly(t *T) {
	subject := UpdatePackages
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestPackages_ListLinkedPackages_Good(t *T) {
	subject := ListLinkedPackages
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestPackages_ListLinkedPackages_Bad(t *T) {
	subject := ListLinkedPackages
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestPackages_ListLinkedPackages_Ugly(t *T) {
	subject := ListLinkedPackages
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}
