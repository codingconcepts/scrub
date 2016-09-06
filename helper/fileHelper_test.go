package helper

import "testing"

func TestGetAllFiles(t *testing.T) {
	fs := setup()
	files, _ := GetAllFiles(fs, []string{"."})

	if len(files) != 6 {
		t.Fatalf("Expected 6 but got %d", len(files))
	}

	assertFile(t, files[0], "oa.txt", 18)
	assertFile(t, files[1], "ob.txt", 18)
	assertFile(t, files[2], "outer_a/ia.txt", 26)
	assertFile(t, files[3], "outer_a/inner_a/iia.txt", 35)
	assertFile(t, files[4], "outer_b/ib.txt", 26)
	assertFile(t, files[5], "outer_b/inner_b/iib.txt", 35)
}

func TestGetTopLevelDirectories(t *testing.T) {
	fs := setup()
	directories, _ := GetTopLevelDirectories(fs, []string{"outer_a", "outer_b", "oa.txt"})

	if len(directories) != 2 {
		t.Fatalf("Expected 2 directories but got %d", len(directories))
	}

	assertDirectory(t, directories[0], "outer_a")
	assertDirectory(t, directories[1], "outer_b")
}
