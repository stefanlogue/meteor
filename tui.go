package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
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
	titleStyle      lipgloss.Style
	headerText      = lipgloss.NewStyle().Foreground(lipgloss.Color("#F48F0B")).Bold(true).Padding(0, 1, 0, 2)
	promptStyle     = lipgloss.NewStyle().Margin(1, 0, 0, 0)
	selectedStyle   = lipgloss.NewStyle().Background(lipgloss.Color("212")).Foreground(lipgloss.Color("230")).Padding(0, 3).Margin(1, 1)
	unselectedStyle = lipgloss.NewStyle().Background(lipgloss.Color("235")).Foreground(lipgloss.Color("254")).Padding(0, 3).Margin(1, 1)
	docStyle        = lipgloss.NewStyle().Margin(1, 2)
	inputStyle      = lipgloss.NewStyle().
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

type coauthorListKeyMap struct {
	selectItem key.Binding
	accept     key.Binding
	exit       key.Binding
}

func newCoauthorListKeyMap() *coauthorListKeyMap {
	return &coauthorListKeyMap{
		selectItem: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "select/deselect"),
		),
		accept: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "accept"),
		),
		exit: key.NewBinding(
			key.WithKeys("ctrl+c", "q"),
			key.WithHelp("ctrl+c/q", "quit"),
		),
	}
}

type commitMessageShortInputKeyMap struct {
	submit key.Binding
	exit   key.Binding
}

func (c commitMessageShortInputKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		c.submit,
		c.exit,
	}
}

func (c commitMessageShortInputKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		c.ShortHelp(),
	}
}

func newCommitMessageShortInputKeyMap() *commitMessageShortInputKeyMap {
	return &commitMessageShortInputKeyMap{
		submit: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "submit"),
		),
		exit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit"),
		),
	}
}

type commitMessageLongInputKeyMap struct {
	submit key.Binding
	exit   key.Binding
}

func (c commitMessageLongInputKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		c.submit,
		c.exit,
	}
}

func (c commitMessageLongInputKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		c.ShortHelp(),
	}
}

func newCommitMessageLongInputKeyMap() *commitMessageLongInputKeyMap {
	return &commitMessageLongInputKeyMap{
		submit: key.NewBinding(
			key.WithKeys("ctrl+d"),
			key.WithHelp("ctrl+d", "submit"),
		),
		exit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit"),
		),
	}
}

/* Model */
type Model struct {
	boardList                   list.Model
	prefixList                  list.Model
	coauthorList                list.Model
	commitMessageLongInput      textarea.Model
	commitMessageShortInput     textinput.Model
	ticketNumberInput           textinput.Model
	commitMessageLongInputHelp  help.Model
	commitMessageShortInputHelp help.Model
	coauthorListKeys            *coauthorListKeyMap
	commitMessageShortInputKeys *commitMessageShortInputKeyMap
	commitMessageLongInputKeys  *commitMessageLongInputKeyMap
	selectedBoard               string
	breakingChangePrompt        string
	coauthorsString             string
	commitMessageShort          string
	commitMessageLong           string
	selectedPrefix              string
	ticketNumber                string
	selectedCoauthorsString     string
	selectedCoauthors           []coauthor
	confirmation                bool
	finished                    bool
	hasBoards                   bool
	hasBreakingChange           bool
	hasCoauthors                bool
	hasCommitMessageShort       bool
	hasCommitMessageLong        bool
	hasCommitted                bool
	hasSelectedBoard            bool
	hasSelectedCoauthors        bool
	hasSelectedPrefix           bool
	hasTicketNumber             bool
	isBreakingChange            bool
	quitting                    bool
	width                       int
}

func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}

func newModel(boards []list.Item, prefixes []list.Item, coauthors []list.Item) *Model {
	var boardList list.Model
	var ticketInput textinput.Model
	hasBoards := false
	hasCoauthors := false
	commitShortKeys := newCommitMessageShortInputKeyMap()
	commitLongKeys := newCommitMessageLongInputKeyMap()

	if boards != nil {
		hasBoards = true
		// set up the boards list
		bDelegate := list.NewDefaultDelegate()
		bDelegate.ShowDescription = false
		boardList = list.New(boards, bDelegate, defaultWidth, listHeight)
		boardList.AdditionalShortHelpKeys = func() []key.Binding {
			return []key.Binding{
				key.NewBinding(
					key.WithKeys("enter"),
					key.WithHelp("enter", "select"),
				),
			}
		}
		boardList.AdditionalFullHelpKeys = func() []key.Binding {
			return []key.Binding{
				key.NewBinding(
					key.WithKeys("enter"),
					key.WithHelp("enter", "select"),
				),
			}
		}
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
	coListKeys := newCoauthorListKeyMap()
	if coauthors != nil {
		hasCoauthors = true
		// set up the co-authors list
		coListDelegate := list.NewDefaultDelegate()
		coListDelegate.ShortHelpFunc = func() []key.Binding {
			return []key.Binding{
				coListKeys.selectItem,
				coListKeys.accept,
				coListKeys.exit,
			}
		}
		colist = list.New(coauthors, coListDelegate, defaultWidth, listHeight)
		colist.Styles.Title = colist.Styles.Title.Copy().Background(lipgloss.Color("#F48F0B"))
		colist.Title = "co-author(s)"
		colist.SetShowStatusBar(false)
		colist.SetFilteringEnabled(true)
	}

	// set up the prefixes list
	prefixList := list.New(prefixes, list.NewDefaultDelegate(), defaultWidth, listHeight)
	prefixList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "select"),
			),
		}
	}
	prefixList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "select"),
			),
		}
	}
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

	return &Model{
		boardList:                   boardList,
		prefixList:                  prefixList,
		coauthorList:                colist,
		ticketNumberInput:           ticketInput,
		commitMessageShortInput:     cms,
		commitMessageLongInput:      cml,
		selectedCoauthorsString:     "selected: ",
		hasBoards:                   hasBoards,
		hasCoauthors:                hasCoauthors,
		isBreakingChange:            false,
		breakingChangePrompt:        "Is this a breaking change? (y/n)",
		coauthorListKeys:            coListKeys,
		commitMessageLongInputHelp:  help.New(),
		commitMessageShortInputHelp: help.New(),
		commitMessageLongInputKeys:  commitLongKeys,
		commitMessageShortInputKeys: commitShortKeys,
		width:                       defaultWidth,
	}
}

func (m *Model) updateBoardList(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m *Model) updatePrefixList(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if coauthor.Name == "None" {
				return ""
			}
			s += fmt.Sprintf("\nCo-authored-by: %s <%s>", coauthor.Name, coauthor.Email)
		}
	}
	return s
}

func removeCoauthor(coauthors []coauthor, coauthor coauthor) []coauthor {
	for i, c := range coauthors {
		if c.Name == coauthor.Name {
			return append(coauthors[:i], coauthors[i+1:]...)
		}
	}
	return coauthors
}

func (m *Model) updateCoauthorList(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.coauthorList.SetSize(msg.Width-h, msg.Height-v)
		return m, nil

	case tea.KeyMsg:
		if m.coauthorList.FilterState() == list.Filtering {
			break
		}
		switch {
		case key.Matches(msg, m.coauthorListKeys.exit):
			m.quitting = true
			return m, tea.Quit

		case key.Matches(msg, m.coauthorListKeys.selectItem):
			i, ok := m.coauthorList.SelectedItem().(coauthor)
			if ok {
				i.Selected = !i.Selected
				if i.Selected {
					if i.Name == "None" {
						m.selectedCoauthors = []coauthor{}
					}
					m.selectedCoauthors = append(m.selectedCoauthors, i)
				} else {
					m.selectedCoauthors = removeCoauthor(m.selectedCoauthors, i)
				}
				index := m.coauthorList.Index()
				m.coauthorList.SetItem(index, i)
				var coauthorNames []string
				for _, coauthor := range m.selectedCoauthors {
					coauthorNames = append(coauthorNames, coauthor.Name)
				}
				m.selectedCoauthorsString = strings.Join(coauthorNames, ", ")
				if len(m.selectedCoauthors) == 0 {
					m.coauthorList.NewStatusMessage(statusMessageStyle(""))
				} else {
					m.coauthorList.NewStatusMessage(statusMessageStyle(m.selectedCoauthorsString))
				}
			}
			return m, nil
		case key.Matches(msg, m.coauthorListKeys.accept):
			m.hasSelectedCoauthors = true
			m.coauthorsString = buildCoauthorsString(m.selectedCoauthors)
			m.commitMessageLong += m.coauthorsString
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.coauthorList, cmd = m.coauthorList.Update(msg)
	return m, cmd
}

func (m *Model) updateTicketNumberInput(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m *Model) updateBreakingChange(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.confirmation = false
			m.quitting = true
			return m, tea.Quit
		case "n", "N":
			m.confirmation = false
			m.isBreakingChange = false
			m.hasBreakingChange = true
			if !m.hasBoards {
				m.commitMessageShortInput.SetValue(fmt.Sprintf("%s: ", m.selectedPrefix))
			} else {
				m.commitMessageShortInput.SetValue(fmt.Sprintf("%s: <%s> ", m.ticketNumber, m.selectedPrefix))
			}
			m.commitMessageShortInput.Focus()
			return m, nil
		case "y", "Y":
			m.confirmation = true
			m.isBreakingChange = true
			m.hasBreakingChange = true
			if !m.hasBoards {
				m.commitMessageShortInput.SetValue(fmt.Sprintf("%s!: ", m.selectedPrefix))
			} else {
				m.commitMessageShortInput.SetValue(fmt.Sprintf("%s: <%s!> ", m.ticketNumber, m.selectedPrefix))
			}
			m.commitMessageShortInput.Focus()
			return m, nil
		case "l", "h", "tab", "shift+tab", "left", "right", "ctrl+p", "ctrl+n":
			m.confirmation = !m.confirmation
		case "enter":
			m.isBreakingChange = m.confirmation
			m.hasBreakingChange = true
			if !m.hasBoards {
				if !m.isBreakingChange {
					m.commitMessageShortInput.SetValue(fmt.Sprintf("%s: ", m.selectedPrefix))
				} else {
					m.commitMessageShortInput.SetValue(fmt.Sprintf("%s!: ", m.selectedPrefix))
				}
			} else {
				if !m.isBreakingChange {
					m.commitMessageShortInput.SetValue(fmt.Sprintf("%s: <%s> ", m.ticketNumber, m.selectedPrefix))
				} else {
					m.commitMessageShortInput.SetValue(fmt.Sprintf("%s: <%s!> ", m.ticketNumber, m.selectedPrefix))
				}
			}
			m.commitMessageShortInput.Focus()
			return m, nil
		}
	}
	return m, nil
}

func (m *Model) updateCommitMessageShortInput(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m *Model) updateCommitMessageLongInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlD:
			m.commitMessageLong = m.commitMessageLongInput.Value()
			m.hasCommitMessageLong = true
			m.finished = true
			m.quitting = true
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitting = true
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.commitMessageLongInput, cmd = m.commitMessageLongInput.Update(msg)
	return m, cmd
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		m.width = msg.Width
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
	case !m.hasBreakingChange:
		return m.updateBreakingChange(msg)
	case m.hasCoauthors && !m.hasSelectedCoauthors:
		return m.updateCoauthorList(msg)
	case !m.hasCommitMessageShort:
		return m.updateCommitMessageShortInput(msg)
	case !m.hasCommitMessageLong:
		return m.updateCommitMessageLongInput(msg)
	default:
		return m, tea.Quit
	}
}

func (m *Model) CommitMessage() (string, string) {
	m.commitMessageLong = m.commitMessageLong + m.coauthorsString
	return m.commitMessageShort, m.commitMessageLong
}

func (m *Model) Finished() bool {
	return m.hasCommitMessageShort
}

func (m *Model) View() string {
	var yes, no string
	if m.confirmation {
		yes = selectedStyle.Render("Yes")
		no = unselectedStyle.Render("No")
	} else {
		yes = unselectedStyle.Render("Yes")
		no = selectedStyle.Render("No")
	}

	commitMsgShortHelpView := m.commitMessageShortInputHelp.View(m.commitMessageShortInputKeys)
	commitMsgLongHelpView := m.commitMessageLongInputHelp.View(m.commitMessageLongInputKeys)

	header := m.appBoundaryView("meteor")

	s := ""
	switch {
	case m.hasBoards && !m.hasSelectedBoard:
		s = lipgloss.JoinVertical(lipgloss.Top, m.boardList.View(), m.selectedBoard)
	case m.hasBoards && !m.hasTicketNumber:
		title := titleStyle.Render("Ticket number")
		s = lipgloss.NewStyle().MarginLeft(2).Render(lipgloss.JoinVertical(lipgloss.Top, title, inputStyle.MarginTop(1).Render(m.ticketNumberInput.View())))
	case !m.hasSelectedPrefix:
		s = lipgloss.JoinVertical(lipgloss.Top, m.prefixList.View())
	case !m.hasBreakingChange:
		s = lipgloss.JoinVertical(lipgloss.Center, promptStyle.Render(m.breakingChangePrompt), lipgloss.JoinHorizontal(lipgloss.Left, no, yes))
	case m.hasCoauthors && !m.hasSelectedCoauthors:
		s = lipgloss.JoinVertical(lipgloss.Top, m.coauthorList.View())
	case !m.hasCommitMessageShort:
		s = lipgloss.JoinVertical(lipgloss.Top, m.commitMessageShortInput.View(), " ", m.commitMessageLongInput.View(), " ", commitMsgShortHelpView)
	case !m.hasCommitMessageLong:
		s = lipgloss.JoinVertical(lipgloss.Top, m.commitMessageShortInput.View(), " ", m.commitMessageLongInput.View(), " ", commitMsgLongHelpView)

	}
	return header + "\n" + s
}

func (m *Model) appBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		headerText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("#F48F0B")),
	)
}
