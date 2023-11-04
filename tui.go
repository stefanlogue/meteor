package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	defaultWidth = 80
	listHeight   = 35
)

var (
	titleStyle lipgloss.Style
	docStyle   = lipgloss.NewStyle().Margin(1, 2)
	inputStyle = lipgloss.NewStyle().
			BorderForeground(lipgloss.Color("#F48F0B")).
			BorderStyle(lipgloss.RoundedBorder()).
			Padding(1)
	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#F48F0B", Dark: "#F48F0B"}).
				Render
)

func (b board) Title() string       { return b.Name }
func (b board) Description() string { return "" }
func (b board) FilterValue() string { return b.Name }

func (p prefix) Title() string       { return p.T }
func (p prefix) Description() string { return p.D }
func (p prefix) FilterValue() string { return p.T }

func (c coauthor) Title() string       { return c.Name }
func (c coauthor) Description() string { return c.Email }
func (c coauthor) FilterValue() string { return c.Name }

/* Model */
type Model struct {
	boardList               list.Model
	prefixList              list.Model
	coauthorList            list.Model
	commitMessageLongInput  textarea.Model
	commitMessageShortInput textinput.Model
	ticketNumberInput       textinput.Model
	selectedBoard           string
	coauthorsString         string
	commitMessageShort      string
	commitMessageLong       string
	selectedPrefix          string
	ticketNumber            string
	selectedCoauthorsString string
	selectedCoauthors       []coauthor
	finished                bool
	hasBoards               bool
	hasCoauthors            bool
	hasCommitMessageShort   bool
	hasCommitMessageLong    bool
	hasCommitted            bool
	hasSelectedBoard        bool
	hasSelectedCoauthors    bool
	hasSelectedPrefix       bool
	hasTicketNumber         bool
	quitting                bool
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func newModel(boards []list.Item, prefixes []list.Item, coauthors []list.Item) Model {
	var boardList list.Model
	var ticketInput textinput.Model
	hasBoards := false
	hasCoauthors := false

	if boards != nil {
		hasBoards = true
		// set up the boards list
		bDelegate := list.NewDefaultDelegate()
		bDelegate.ShowDescription = false
		boardList = list.New(boards, bDelegate, defaultWidth, listHeight)
		boardList.Styles.Title = boardList.Styles.Title.Copy().Background(lipgloss.Color("#F48F0B"))
		titleStyle = boardList.Styles.Title.Copy()
		boardList.Title = "Select the JIRA board"
		boardList.SetShowStatusBar(false)
		boardList.SetFilteringEnabled(true)

		// set up the ticket number input
		ticketInput = textinput.New()
		ticketInput.Placeholder = "Ticket number"
		ticketInput.CharLimit = 10
	}

	var colist list.Model
	if coauthors != nil {
		hasCoauthors = true
		// set up the co-authors list
		coListDelegate := list.NewDefaultDelegate()
		colist = list.New(coauthors, coListDelegate, defaultWidth, listHeight)
		colist.Styles.Title = colist.Styles.Title.Copy().Background(lipgloss.Color("#F48F0B"))
		colist.Title = "co-author(s)"
		colist.SetShowStatusBar(false)
		colist.SetFilteringEnabled(true)
	}

	// set up the prefixes list
	prefixList := list.New(prefixes, list.NewDefaultDelegate(), defaultWidth, listHeight)
	prefixList.Styles.Title = prefixList.Styles.Title.Copy().Background(lipgloss.Color("#F48F0B"))
	prefixList.Title = "Select the type of change you're committing"
	prefixList.SetShowStatusBar(false)
	prefixList.SetFilteringEnabled(true)

	// set up the short commit message input
	cms := textinput.New()
	cms.Placeholder = "short commit message"
	cms.CharLimit = 72

	// set up the long commit message input
	cml := textarea.New()
	cml.Placeholder = "longer commit message"

	return Model{
		boardList:               boardList,
		prefixList:              prefixList,
		coauthorList:            colist,
		ticketNumberInput:       ticketInput,
		commitMessageShortInput: cms,
		commitMessageLongInput:  cml,
		selectedCoauthorsString: "selected: ",
		hasBoards:               hasBoards,
		hasCoauthors:            hasCoauthors,
	}
}

func (m Model) updateBoardList(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.boardList.SetSize(msg.Width-h, msg.Height-v)
		return m, nil

	case tea.KeyMsg:
		if m.boardList.FilterState() == list.Filtering {
			break
		}
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "enter", " ":
			var ticketPrefix string
			i, ok := m.boardList.SelectedItem().(board)
			if ok {
				m.selectedBoard = i.Name
				m.hasSelectedBoard = true
				if m.selectedBoard == "NONE" {
					m.ticketNumber = "noticket"
					m.hasTicketNumber = true
				} else {
					ticketPrefix = getGitTicketNumber(m.selectedBoard)
				}
				m.ticketNumberInput.Focus()
				if ticketPrefix == "" {
					ticketPrefix = m.selectedBoard + "-"
				}
				m.ticketNumberInput.SetValue(ticketPrefix)
			}
		}
	}

	var cmd tea.Cmd
	m.boardList, cmd = m.boardList.Update(msg)
	return m, cmd
}

func (m Model) updatePrefixList(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.prefixList.SetSize(msg.Width-h, msg.Height-v)
		return m, nil

	case tea.KeyMsg:
		if m.prefixList.FilterState() == list.Filtering {
			break
		}
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "enter", " ":
			i, ok := m.prefixList.SelectedItem().(prefix)
			if ok {
				m.selectedPrefix = i.T
				m.hasSelectedPrefix = true
				if !m.hasBoards {
					m.commitMessageShortInput.SetValue(fmt.Sprintf("%s: ", m.selectedPrefix))
				} else {
					m.commitMessageShortInput.SetValue(fmt.Sprintf("%s: <%s> ", m.ticketNumber, m.selectedPrefix))
				}
				m.commitMessageShortInput.Focus()
			}
		}
	}

	var cmd tea.Cmd
	m.prefixList, cmd = m.prefixList.Update(msg)
	return m, cmd
}

func buildCoauthorsString(coauthors []coauthor) string {
	s := `


	`
	for _, coauthor := range coauthors {
		if coauthor.Selected {
			s += fmt.Sprintf("\nCo-authored-by: %s <%s>", coauthor.Name, coauthor.Email)
		}
	}
	return s
}

func (m Model) updateCoauthorList(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.coauthorList.SetSize(msg.Width-h, msg.Height-v)
		return m, nil

	case tea.KeyMsg:
		if m.coauthorList.FilterState() == list.Filtering {
			break
		}
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case " ":
			i, ok := m.coauthorList.SelectedItem().(coauthor)
			if ok {
				i.Selected = !i.Selected
				m.selectedCoauthors = append(m.selectedCoauthors, i)
				if len(m.selectedCoauthors) == 1 {
					m.selectedCoauthorsString += i.Name
				} else {
					m.selectedCoauthorsString += ", " + i.Name
				}
				m.coauthorList.NewStatusMessage(statusMessageStyle(m.selectedCoauthorsString))
			}
		case "enter":
			m.hasSelectedCoauthors = true
			m.coauthorsString = buildCoauthorsString(m.selectedCoauthors)
			m.commitMessageLong += m.coauthorsString
			m.finished = true
		}
	}
	var cmd tea.Cmd
	m.coauthorList, cmd = m.coauthorList.Update(msg)
	return m, cmd
}

func (m Model) updateTicketNumberInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.ticketNumber = m.ticketNumberInput.Value()
			m.hasTicketNumber = true
			return m, nil
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitting = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.ticketNumberInput, cmd = m.ticketNumberInput.Update(msg)
	return m, cmd
}

func (m Model) updateCommitMessageShortInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.commitMessageShort = m.commitMessageShortInput.Value()
			m.hasCommitMessageShort = true
			m.commitMessageLongInput.Focus()
			return m, nil
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitting = true
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.commitMessageShortInput, cmd = m.commitMessageShortInput.Update(msg)
	return m, cmd
}

func (m Model) updateCommitMessageLongInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlD, tea.KeyEsc:
			m.commitMessageLong = m.commitMessageLongInput.Value()
			m.hasCommitMessageLong = true
			if !m.hasCoauthors {
				m.finished = true
			}
			return m, nil
		case tea.KeyCtrlC:
			m.quitting = true
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.commitMessageLongInput, cmd = m.commitMessageLongInput.Update(msg)
	return m, cmd
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		if m.hasBoards {
			m.boardList.SetSize(msg.Width-h, msg.Height-v)
		}
		m.prefixList.SetSize(msg.Width-h, msg.Height-v)
		if m.hasCoauthors {
			m.coauthorList.SetSize(msg.Width-h, msg.Height-v)
		}
		m.commitMessageShortInput.Width = msg.Width - h - 4
		m.commitMessageLongInput.SetWidth(msg.Width - h - 4)
		m.commitMessageLongInput.SetHeight(msg.Height - v - 4)
		inputStyle = inputStyle.Width(msg.Width - h - 4)
		return m, nil
	}
	switch {
	case m.hasBoards && !m.hasSelectedBoard:
		return m.updateBoardList(msg)
	case m.hasBoards && !m.hasTicketNumber:
		return m.updateTicketNumberInput(msg)
	case !m.hasSelectedPrefix:
		return m.updatePrefixList(msg)
	case !m.hasCommitMessageShort:
		return m.updateCommitMessageShortInput(msg)
	case !m.hasCommitMessageLong:
		return m.updateCommitMessageLongInput(msg)
	case m.hasCoauthors && !m.hasSelectedCoauthors:
		return m.updateCoauthorList(msg)
	default:
		return m, tea.Quit
	}
}

func (m Model) CommitMessage() (string, string) {
	return m.commitMessageShort, m.commitMessageLong
}

func (m Model) Finished() bool {
	return m.finished
}

func (m Model) View() string {
	s := ""
	switch {
	case m.hasBoards && !m.hasSelectedBoard:
		s = lipgloss.JoinVertical(lipgloss.Top, m.boardList.View(), m.selectedBoard)
	case m.hasBoards && !m.hasTicketNumber:
		title := titleStyle.Render("Ticket number")
		s = lipgloss.NewStyle().MarginLeft(2).Render(lipgloss.JoinVertical(lipgloss.Top, title, inputStyle.MarginTop(1).Render(m.ticketNumberInput.View())))
	case !m.hasSelectedPrefix:
		s = lipgloss.JoinVertical(lipgloss.Top, m.prefixList.View())
	case !m.hasCommitMessageShort, !m.hasCommitMessageLong:
		s = lipgloss.JoinVertical(lipgloss.Top, m.commitMessageShortInput.View(), " ", m.commitMessageLongInput.View())
	case m.hasCoauthors && !m.hasSelectedCoauthors:
		s = lipgloss.JoinVertical(lipgloss.Top, m.coauthorList.View())
	case m.quitting:
		s = "Goodbye!"
	}
	return s
}
