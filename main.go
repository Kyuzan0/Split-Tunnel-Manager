package main

import (
	"image/color"
	"log"
	"os"
	"path/filepath"
	"syscall"

	"golang.org/x/sys/windows"

	"fyne.io/fyne/v2"
	fyneapp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"

	"split-tunnel-manager/internal/netcheck"
	"split-tunnel-manager/internal/split"
	"split-tunnel-manager/internal/store"
	"split-tunnel-manager/internal/ui/shell"
)

func isAdmin() bool {
	var sid *windows.SID
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid)
	if err != nil {
		return false
	}
	defer windows.FreeSid(sid)

	token := windows.Token(0)
	member, err := token.IsMember(sid)
	if err != nil {
		return false
	}
	return member
}

func elevate() {
	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	verb := syscall.StringToUTF16Ptr("runas")
	exePtr := syscall.StringToUTF16Ptr(exe)
	cwdPtr := syscall.StringToUTF16Ptr(cwd)
	var argString string
	for _, arg := range os.Args[1:] {
		argString += arg + " "
	}
	argPtr := syscall.StringToUTF16Ptr(argString)

	err = windows.ShellExecute(0, verb, exePtr, argPtr, cwdPtr, syscall.SW_SHOW)
	if err != nil {
		log.Fatal("Gagal meminta izin Administrator:", err)
	}
	os.Exit(0)
}

func main() {
	if !isAdmin() {
		elevate()
	}

	dataDir, err := appDataDir()
	if err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(dataDir, 0o700); err != nil {
		log.Fatal(err)
	}

	st, err := store.NewStore(filepath.Join(dataDir, "split-tunnel.db"))
	if err != nil {
		log.Fatal(err)
	}
	defer st.Close()

	spl := split.NewManager()
	chk := netcheck.NewChecker("", 0)

	a := fyneapp.NewWithID("com.splittunnel.manager")
	a.Settings().SetTheme(&customTheme{})
	w := a.NewWindow("Split Tunnel Manager")
	w.Resize(fyne.NewSize(750, 550))
	w.SetContent(shell.New(shell.Deps{
		Store: st,
		Split: spl,
		Net:   chk,
	}))
	w.ShowAndRun()
}

// customTheme menerapkan tampilan dark/modern dengan warna aksen biru-ungu.
type customTheme struct{}

func (c *customTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	switch n {
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 0x5D, G: 0x3E, B: 0xDE, A: 0xFF} // Deep purple/blue
	case theme.ColorNameButton:
		if v == theme.VariantLight {
			return color.NRGBA{R: 0xE0, G: 0xE0, B: 0xE0, A: 0xFF}
		}
		return color.NRGBA{R: 0x2C, G: 0x2C, B: 0x2E, A: 0xFF} // Darker button
	case theme.ColorNameBackground:
		if v == theme.VariantLight {
			return color.NRGBA{R: 0xF5, G: 0xF5, B: 0xF5, A: 0xFF}
		}
		return color.NRGBA{R: 0x1A, G: 0x1A, B: 0x1C, A: 0xFF} // Very dark background
	}
	// Fallback ke default
	if v == theme.VariantLight {
		return theme.DefaultTheme().Color(n, theme.VariantLight)
	}
	return theme.DefaultTheme().Color(n, theme.VariantDark)
}

func (c *customTheme) Font(s fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(s)
}

func (c *customTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (c *customTheme) Size(n fyne.ThemeSizeName) float32 {
	if n == theme.SizeNamePadding {
		return 8.0 // Padding lebih lebar
	}
	if n == theme.SizeNameText {
		return 14.0 // Teks sedikit lebih besar
	}
	return theme.DefaultTheme().Size(n)
}

func appDataDir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "Split Tunnel Manager"), nil
}
