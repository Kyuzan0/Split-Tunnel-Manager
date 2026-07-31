package shell

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"

	"split-tunnel-manager/internal/domain"
	"split-tunnel-manager/internal/netcheck"
	"split-tunnel-manager/internal/split"
	"split-tunnel-manager/internal/store"
)

// Deps wires services into the Fyne shell
type Deps struct {
	Store *store.Store
	Split *split.Manager
	Net   *netcheck.Checker
}

func New(d Deps) fyne.CanvasObject {
	s := &shell{deps: d}
	s.connectTab = s.buildConnect()
	s.settingsTab = s.buildSettings()

	tabs := container.NewAppTabs(
		container.NewTabItem("Split Tunnel", s.connectTab),
		container.NewTabItem("Bypass Rules", s.settingsTab),
	)
	tabs.SetTabLocation(container.TabLocationLeading)
	tabs.OnSelected = func(ti *container.TabItem) {
		switch ti.Text {
		case "Split Tunnel":
			s.refreshConnect()
		case "Bypass Rules":
			s.refreshBypass()
		}
	}
	s.refreshBypass()
	s.refreshConnect()
	return tabs
}

type shell struct {
	deps Deps

	connectTab, settingsTab fyne.CanvasObject

	// Bypass State
	bypassList     *widget.List
	bypassRules    []domain.BypassRule
	selectedBypass string
	cidrEntry      *widget.Entry
	cidrLabelEnt   *widget.Entry

	// Connect State
	adapterSelect *widget.Select
	statusLabel   *widget.Label
	splitStatus   domain.SplitStatus
}

func (s *shell) buildConnect() fyne.CanvasObject {
	s.statusLabel = widget.NewLabel("")
	s.statusLabel.Wrapping = fyne.TextWrapWord
	s.adapterSelect = widget.NewSelect([]string{"Memuat interface..."}, nil)
	s.adapterSelect.SetSelected("Memuat interface...")

	loadAdapters := func() {
		adapters, err := netcheck.GetActiveAdapters()
		if err != nil {
			dialog.ShowError(err, primaryWindow())
			return
		}
		var names []string
		for _, a := range adapters {
			names = append(names, fmt.Sprintf("%s - %s", a.Name, a.Desc))
		}
		
		s.adapterSelect.Options = names
		if len(names) > 0 {
			s.adapterSelect.SetSelected(names[0])
		} else {
			s.adapterSelect.Options = []string{"Tidak ada adapter"}
			s.adapterSelect.SetSelected("Tidak ada adapter")
		}
	}

	refreshAdaptersBtn := widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), loadAdapters)

	// Auto-load on startup in background
	go loadAdapters()

	var toggleBtn *widget.Button
	toggleBtn = widget.NewButtonWithIcon("Start Split Tunnel", theme.MediaPlayIcon(), func() {
		win := primaryWindow()
		
		// Jika status saat ini adalah CONNECTED, maka aksi selanjutnya adalah STOP
		if s.splitStatus.State == domain.ConnConnected {
			_ = s.deps.Split.RemoveAll()
			s.splitStatus.State = domain.ConnIdle
			s.splitStatus.InterfaceName = ""
			
			toggleBtn.SetText("Start Split Tunnel")
			toggleBtn.SetIcon(theme.MediaPlayIcon())
			toggleBtn.Importance = widget.HighImportance
			toggleBtn.Refresh()
			
			dialog.ShowInformation("Info", "Split Tunnel telah dimatikan (Bypass dicabut).", win)
			s.refreshConnect()
			return
		}

		// Jika status saat ini idle/error, maka aksi selanjutnya adalah START
		if s.adapterSelect.Selected == "" || s.adapterSelect.Selected == "Tidak ada adapter" || s.adapterSelect.Selected == "Memuat interface..." {
			dialog.ShowInformation("Error", "Pilih Network Interface VPN (misal: TAP-Windows / WireGuard) terlebih dahulu.", win)
			return
		}

		activeBypass := enabledRules(s.bypassRules)
		if len(activeBypass) == 0 {
			dialog.ShowInformation("Info", "Tidak ada rule bypass yang aktif.", win)
			return
		}

		s.statusLabel.SetText("Applying routes...")
		
		if err := s.deps.Split.Apply(activeBypass); err != nil {
			dialog.ShowError(fmt.Errorf("gagal apply split tunnel: %v\n\nPastikan run sebagai Administrator!", err), win)
			s.splitStatus.State = domain.ConnError
			s.splitStatus.LastError = err.Error()
		} else {
			s.splitStatus.State = domain.ConnConnected
			s.splitStatus.InterfaceName = s.adapterSelect.Selected
			s.splitStatus.LastError = ""

			toggleBtn.SetText("Stop Split Tunnel")
			toggleBtn.SetIcon(theme.MediaStopIcon())
			toggleBtn.Importance = widget.DangerImportance
			toggleBtn.Refresh()
			
			dialog.ShowInformation("Success", "Split tunnel routing berhasil diterapkan.", win)
		}
		
		s.refreshConnect()
	})
	toggleBtn.Importance = widget.HighImportance

	form := widget.NewForm(
		widget.NewFormItem("VPN Interface", s.adapterSelect),
	)
	
	card := widget.NewCard("Split Tunnel Manager", "Pastikan VPN Anda (Avira/HMA) SUDAH CONNECT sebelum Start.", form)

	return container.NewBorder(
		card,
		container.NewHBox(refreshAdaptersBtn, toggleBtn),
		nil, nil,
		container.NewVBox(s.statusLabel),
	)
}

func (s *shell) refreshConnect() {
	if s.statusLabel == nil {
		return
	}
	var b strings.Builder
	
	b.WriteString(fmt.Sprintf("Status: %s\n", s.splitStatus.State))
	if s.splitStatus.InterfaceName != "" {
		b.WriteString(fmt.Sprintf("Applied on: %s\n", s.splitStatus.InterfaceName))
	}
	if s.splitStatus.LastError != "" {
		b.WriteString(fmt.Sprintf("Error: %s\n", s.splitStatus.LastError))
	}

	active := enabledRules(s.bypassRules)
	if len(active) == 0 {
		b.WriteString("\nActive bypass: (none)\n")
	} else {
		b.WriteString("\nActive bypass rules to route to Gateway:\n")
		for _, r := range active {
			b.WriteString(fmt.Sprintf("  • %s (%s)\n", r.CIDR, r.Label))
		}
	}
	s.statusLabel.SetText(b.String())
}

func (s *shell) buildSettings() fyne.CanvasObject {
	s.bypassList = widget.NewList(
		func() int { return len(s.bypassRules) },
		func() fyne.CanvasObject { 
			icon := widget.NewIcon(theme.FolderIcon())
			label := widget.NewLabel("cidr")
			return container.NewHBox(icon, label)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			if i < 0 || i >= len(s.bypassRules) {
				return
			}
			r := s.bypassRules[i]
			en := "off"
			icon := theme.ErrorIcon()
			if r.Enabled {
				en = "on"
				icon = theme.ConfirmIcon()
			}
			
			box := o.(*fyne.Container)
			box.Objects[0].(*widget.Icon).SetResource(icon)
			box.Objects[1].(*widget.Label).SetText(fmt.Sprintf("[%s] %s — %s", en, r.CIDR, r.Label))
		},
	)
	s.bypassList.OnSelected = func(id widget.ListItemID) {
		if id >= 0 && id < len(s.bypassRules) {
			s.selectedBypass = s.bypassRules[id].ID
		}
	}

	s.cidrEntry = widget.NewEntry()
	s.cidrEntry.SetPlaceHolder("192.168.2.0/24")
	s.cidrLabelEnt = widget.NewEntry()
	s.cidrLabelEnt.SetPlaceHolder("Label LAN")

	addBtn := widget.NewButtonWithIcon("Add CIDR", theme.ContentAddIcon(), func() {
		win := primaryWindow()
		cidr := strings.TrimSpace(s.cidrEntry.Text)
		if cidr == "" {
			return
		}
		label := strings.TrimSpace(s.cidrLabelEnt.Text)
		if label == "" {
			label = cidr
		}
		rule := domain.BypassRule{
			ID:      uuid.NewString(),
			CIDR:    cidr,
			Label:   label,
			Enabled: true,
			Source:  "user",
		}
		
		if err := s.deps.Store.UpsertBypass(rule); err != nil {
			dialog.ShowError(err, win)
			return
		}
		s.cidrEntry.SetText("")
		s.cidrLabelEnt.SetText("")
		s.refreshBypass()
		s.refreshConnect()
	})
	addBtn.Importance = widget.HighImportance

	delBtn := widget.NewButtonWithIcon("Delete selected", theme.DeleteIcon(), func() {
		if s.selectedBypass == "" {
			return
		}
		
		var deletedCIDR string
		for _, r := range s.bypassRules {
			if r.ID == s.selectedBypass {
				deletedCIDR = r.CIDR
				break
			}
		}

		if err := s.deps.Store.DeleteBypass(s.selectedBypass); err != nil {
			dialog.ShowError(err, primaryWindow())
			return
		}
		
		if deletedCIDR != "" {
			_ = s.deps.Split.RemoveRule(deletedCIDR)
		}

		s.selectedBypass = ""
		s.refreshBypass()
		s.refreshConnect()
	})
	delBtn.Importance = widget.DangerImportance

	form := container.NewVBox(
		widget.NewLabel("Bypass CIDR (IP yang TIDAK AKAN masuk VPN)"),
		s.cidrEntry,
		s.cidrLabelEnt,
		container.NewHBox(addBtn, delBtn),
		widget.NewLabel("Modifikasi tabel routing otomatis via UAC."),
	)
	
	card := widget.NewCard("Konfigurasi Aturan Bypass", "", form)

	return container.NewBorder(card, nil, nil, nil, s.bypassList)
}

func (s *shell) refreshBypass() {
	rules, err := s.deps.Store.ListBypass()
	if err != nil {
		s.bypassRules = nil
	} else {
		s.bypassRules = rules
	}
	if s.bypassList != nil {
		s.bypassList.Refresh()
	}
}

func enabledRules(all []domain.BypassRule) []domain.BypassRule {
	out := make([]domain.BypassRule, 0, len(all))
	for _, r := range all {
		if r.Enabled {
			out = append(out, r)
		}
	}
	return out
}

func primaryWindow() fyne.Window {
	wins := fyne.CurrentApp().Driver().AllWindows()
	if len(wins) == 0 {
		return nil
	}
	return wins[0]
}
