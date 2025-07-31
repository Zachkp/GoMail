package models

import (
	"github.com/Zachkp/GoMail/styles"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

func CreateColumns(width int) []table.Column {
	// Reserve a little padding (e.g., for borders or spacing)
	totalWidth := width - 10

	// Allocate width proportionally
	senderWidth := 25
	dateWidth := 10
	timeWidth := 10
	messageWidth := totalWidth - senderWidth - dateWidth - timeWidth

	// Don't let message width go negative
	if messageWidth < 20 {
		messageWidth = 20
	}

	return []table.Column{
		{Title: "Sender", Width: senderWidth},
		{Title: "Date", Width: dateWidth},
		{Title: "Time", Width: timeWidth},
		{Title: "Message", Width: messageWidth},
	}
}

func CreateTable() model {
	columns := CreateColumns(styles.PlaceholderWidth)

	//TODO: Populate from email account, then use Bat to display whole email
	// Just for testing purposes while TUI is developed
	rows := []table.Row{
		{"alice@company.com", "11/15/2024", "9:15AM", "Meeting rescheduled to 3PM today. Please confirm your attendance."},
		{"bob.wilson@startup.io", "11/14/2024", "2:45PM", "Your invoice #12345 is now due. Payment can be made via our online portal."},
		{"newsletter@techblog.com", "11/14/2024", "6:00AM", "Weekly Digest: Latest trends in AI, blockchain updates, and developer tools you shouldn't miss."},
		{"sarah.j@marketing.co", "11/13/2024", "4:20PM", "Campaign results are in! Click-through rates exceeded expectations by 25%."},
		{"support@cloudservice.net", "11/13/2024", "11:30AM", "System maintenance scheduled for tonight 11PM-2AM EST. Minimal downtime expected."},
		{"mom@family.net", "11/12/2024", "7:45PM", "Don't forget dinner this Sunday at 6PM. Bring dessert if you can!"},
		{"recruiter@bigtech.com", "11/12/2024", "1:15PM", "Exciting senior developer position available. Are you interested in learning more?"},
		{"orders@retailstore.com", "11/11/2024", "3:30PM", "Your order #ORD-789 has shipped! Track your package with code TR-456-ABC."},
		{"team@projectalpha.dev", "11/11/2024", "10:00AM", "Sprint review meeting tomorrow at 2PM. Please prepare your status updates."},
		{"alerts@security.org", "11/10/2024", "8:20AM", "Security alert: Suspicious login attempt detected from unknown location."},
		{"jane.doe@consulting.biz", "11/10/2024", "5:50PM", "Project proposal attached. Looking forward to your feedback by Friday."},
		{"events@conference.tech", "11/09/2024", "12:15PM", "Early bird registration ends tomorrow! Save 40% on TechConf 2025 tickets."},
		{"billing@utilities.gov", "11/09/2024", "9:40AM", "Your monthly statement is ready. View online or download PDF attachment."},
		{"updates@news.daily", "11/08/2024", "6:30AM", "Breaking: Major breakthrough in quantum computing announced by research team."},
		{"coach@fitness.club", "11/08/2024", "7:00PM", "Great workout today! Your next personal training session is scheduled for Monday."},
		{"admin@university.edu", "11/07/2024", "2:00PM", "Course registration opens next week. Priority given to seniors and graduate students."},
		{"noreply@bank.secure", "11/07/2024", "8:15AM", "Your account balance is low. Consider setting up automatic transfers to avoid fees."},
		{"travel@airline.fly", "11/06/2024", "4:45PM", "Flight delayed by 2 hours due to weather. We'll keep you updated on changes."},
		{"hello@startup.new", "11/06/2024", "11:20AM", "Product launch party invitation! Join us next Friday for food, drinks, and demos."},
		{"legal@lawfirm.professional", "11/05/2024", "3:15PM", "Contract review completed. Please see attached redlined version for your approval."},
		{"doctor@clinic.health", "11/05/2024", "10:30AM", "Appointment reminder: Annual checkup scheduled for Thursday at 2:30PM."},
		{"pets@veterinary.care", "11/04/2024", "1:45PM", "Vaccination reminder for Fluffy. Please schedule within the next 2 weeks."},
		{"info@library.public", "11/04/2024", "4:10PM", "Reserved book 'The Pragmatic Programmer' is now available for pickup."},
		{"sales@software.saas", "11/03/2024", "9:25AM", "Limited time offer: 50% off annual subscription. Offer expires this weekend."},
		{"garden@community.local", "11/03/2024", "6:15PM", "Community garden workday this Saturday 9AM. Bring gloves and water bottle."},
		{"chef@restaurant.fine", "11/02/2024", "2:30PM", "Special tasting menu available tonight featuring seasonal ingredients."},
		{"mechanic@autorepair.pro", "11/02/2024", "11:45AM", "Your car is ready for pickup. Total came to $287 as estimated."},
		{"teacher@school.k12", "11/01/2024", "3:20PM", "Parent-teacher conferences next week. Please sign up for your preferred time slot."},
		{"weather@alerts.gov", "11/01/2024", "5:40AM", "Severe weather warning issued for your area. Stay indoors if possible."},
		{"editor@magazine.monthly", "10/31/2024", "4:55PM", "Article submission deadline extended to November 15th. Happy Halloween!"},
		{"landlord@property.mgmt", "10/31/2024", "12:00PM", "Rent reminder: Payment due tomorrow. Late fees apply after the 5th."},
		{"volunteer@charity.org", "10/30/2024", "7:30PM", "Thank you for volunteering! Next food drive is scheduled for December 10th."},
		{"coach@sports.league", "10/30/2024", "6:45PM", "Practice cancelled due to field conditions. Make-up session Sunday 10AM."},
		{"accountant@tax.services", "10/29/2024", "1:10PM", "Tax documents ready for review. Please schedule appointment at your convenience."},
		{"friends@socialclub.fun", "10/29/2024", "8:00PM", "Game night this Friday! Bring your favorite board game and snacks to share."},
		{"delivery@logistics.fast", "10/28/2024", "3:25PM", "Package delivery attempted. You weren't home - rescheduled for tomorrow."},
		{"principal@highschool.edu", "10/28/2024", "10:15AM", "School fundraiser exceeded goals! Thank you to all families who participated."},
		{"curator@museum.art", "10/27/2024", "2:50PM", "New exhibit opening next month: 'Digital Art in the Modern Age' - preview night invitation."},
		{"coordinator@workshop.skill", "10/27/2024", "11:35AM", "Photography workshop this weekend still has openings. Register before Thursday."},
		{"manager@retail.chain", "10/26/2024", "4:15PM", "Employee schedule for next week posted. Check the staff portal for updates."},
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(20),
	)

	m := model{
		table: t,
		width: styles.PlaceholderWidth,
	}

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(styles.Green)).
		BorderBottom(true).
		Bold(true)

	s.Selected = s.Selected.
		Foreground(lipgloss.Color(styles.White)).
		Background(lipgloss.Color(styles.DarkGray)).
		Bold(true)

	t.SetStyles(s)

	return m
}
