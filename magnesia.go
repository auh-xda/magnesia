package main

func (Magnesia) Installed() bool {
	_, err := (Config{}).Parse()

	return nil == err
}
