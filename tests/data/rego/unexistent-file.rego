package bark

default allow := false

allow {
    file.exists("./this-is-an-unexistent-file.json")
}