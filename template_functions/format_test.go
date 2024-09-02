package template_functions

import "testing"


func TestDocStringReplaceAllLeadingTabs(t *testing.T) {
	//Test that it only replaces leading tabs
	have := DocString("This is a test\n\twith a non-leading tab")
	want := "This is a test\n\twith a non-leading tab"
	if have != want {
		t.Errorf("expected %q but got %q", want, have)
	}
	// //Test replacing single leading tab
	// have = DocString("\tThis is a test\n\twith a leading tab")
	// want = "    This is a test\nwith a leading tab"
	// if have != want {
	// 	t.Errorf("expected %q but got %q", want, have)
	// }
	// //Test replacing multiple leading tabs
	// have = DocString("\t\t\tThis is a test\n\twith a leading tab")
	// want = "            This is a test\nwith a leading tab"
	// if have != want {
	// 	t.Errorf("expected %q but got %q", want, have)
	// }
}