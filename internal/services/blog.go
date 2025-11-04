package services

import (
	"errors"
	"strings"
	"template-store/internal/models"
	"time"

	"github.com/gomarkdown/markdown"
	"gorm.io/gorm"
)

type BlogService struct {
	db *gorm.DB
}

func NewBlogService(db *gorm.DB) *BlogService {
	return &BlogService{db: db}
}

// CreateBlogPost creates a new blog post
func (s *BlogService) CreateBlogPost(post *models.BlogPost) error {
	if post.Title == "" {
		return errors.New("blog post title is required")
	}
	if post.Content == "" {
		return errors.New("blog post content is required")
	}
	if post.AuthorID == 0 {
		return errors.New("author ID is required")
	}
	
	// Set default values
	if post.CreatedAt.IsZero() {
		post.CreatedAt = time.Now()
	}
	if post.UpdatedAt.IsZero() {
		post.UpdatedAt = time.Now()
	}
	
	return s.db.Create(post).Error
}

// GetBlogPost retrieves a blog post by ID
func (s *BlogService) GetBlogPost(id uint) (*models.BlogPost, error) {
	var post models.BlogPost
	err := s.db.Preload("Author").Preload("Category").First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// ListBlogPosts retrieves all blog posts with optional filtering
func (s *BlogService) ListBlogPosts(categoryID *uint, search *string, limit, offset int) ([]models.BlogPost, int64, error) {
	var posts []models.BlogPost
	var total int64
	
	query := s.db.Model(&models.BlogPost{}).Preload("Author").Preload("Category")
	
	// Apply category filter
	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}
	
	// Apply search filter
	if search != nil && *search != "" {
		query = query.Where("title ILIKE ? OR content ILIKE ?", "%"+*search+"%", "%"+*search+"%")
	}
	
	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Apply pagination and ordering (newest first)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	
	err := query.Order("created_at DESC").Find(&posts).Error
	return posts, total, err
}

// UpdateBlogPost updates an existing blog post
func (s *BlogService) UpdateBlogPost(id uint, updates map[string]interface{}) error {
	if title, exists := updates["title"]; exists && title == "" {
		return errors.New("blog post title cannot be empty")
	}
	if content, exists := updates["content"]; exists && content == "" {
		return errors.New("blog post content cannot be empty")
	}
	
	// Update the updated_at timestamp
	updates["updated_at"] = time.Now()
	
	return s.db.Model(&models.BlogPost{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteBlogPost deletes a blog post by ID
func (s *BlogService) DeleteBlogPost(id uint) error {
	return s.db.Delete(&models.BlogPost{}, id).Error
}

// GetBlogPostsByCategory retrieves blog posts filtered by category
func (s *BlogService) GetBlogPostsByCategory(categoryID uint) ([]models.BlogPost, error) {
	var posts []models.BlogPost
	err := s.db.Where("category_id = ?", categoryID).Preload("Author").Preload("Category").Order("created_at DESC").Find(&posts).Error
	return posts, err
}

// GetBlogPostsByAuthor retrieves blog posts filtered by author
func (s *BlogService) GetBlogPostsByAuthor(authorID uint) ([]models.BlogPost, error) {
	var posts []models.BlogPost
	err := s.db.Where("author_id = ?", authorID).Preload("Author").Preload("Category").Order("created_at DESC").Find(&posts).Error
	return posts, err
}

// ProcessMarkdown converts markdown content to HTML
func (s *BlogService) ProcessMarkdown(content string) string {
	md := []byte(content)
	html := markdown.ToHTML(md, nil, nil)
	return string(html)
}

// GenerateExcerpt generates a short excerpt from blog content
func (s *BlogService) GenerateExcerpt(content string, maxLength int) string {
	// Remove markdown formatting for excerpt
	plainText := s.stripMarkdown(content)
	
	if len(plainText) <= maxLength {
		return plainText
	}
	
	// Truncate and add ellipsis
	excerpt := plainText[:maxLength]
	lastSpace := strings.LastIndex(excerpt, " ")
	if lastSpace > 0 {
		excerpt = excerpt[:lastSpace]
	}
	
	return excerpt + "..."
}

// stripMarkdown removes markdown formatting from text
func (s *BlogService) stripMarkdown(content string) string {
	// Simple markdown stripping - in production, you might want a more robust solution
	content = strings.ReplaceAll(content, "#", "")
	content = strings.ReplaceAll(content, "*", "")
	content = strings.ReplaceAll(content, "**", "")
	content = strings.ReplaceAll(content, "`", "")
	content = strings.ReplaceAll(content, "```", "")
	content = strings.ReplaceAll(content, ">", "")
	content = strings.ReplaceAll(content, "-", "")
	content = strings.ReplaceAll(content, "+", "")
	
	// Remove extra whitespace
	content = strings.Join(strings.Fields(content), " ")
	
	return content
}

// SeedBlogPosts seeds the database with initial blog posts
func (s *BlogService) SeedBlogPosts() error {
	posts := []models.BlogPost{
		{
			Title:      "Understanding Security Governance: Building a Strong Foundation",
			Content:    `# Understanding Security Governance: Building a Strong Foundation

Security governance is the framework that ensures your organization's security policies align with business objectives and regulatory requirements. Whether you're a startup or an established enterprise, having a solid governance structure is essential for protecting your assets and maintaining stakeholder trust.

## What is Security Governance?

Security governance encompasses the policies, procedures, and controls that guide how your organization manages information security. It's not just about technology—it's about creating a culture of security awareness and accountability at every level of your organization.

Think of security governance as the blueprint for your security program. It defines:

- **Who** is responsible for security decisions
- **What** security standards and policies apply
- **How** security measures are implemented and monitored
- **Why** certain controls are necessary

## Key Components of Effective Security Governance

### 1. Clear Policies and Standards

Your security policies should be documented, accessible, and regularly updated. They serve as the rulebook for how your organization handles sensitive information and responds to security incidents.

**Need help getting started?** Check out our [Security Policy Template](/templates/security-policy-framework) to establish comprehensive policies quickly.

### 2. Defined Roles and Responsibilities

Security isn't just the IT department's job. Effective governance requires clear ownership across the organization. This includes:

- Executive leadership providing strategic direction
- Security officers implementing and monitoring controls
- Department managers ensuring team compliance
- All employees following security best practices

Our [RACI Matrix Template](/templates/security-raci-matrix) can help you define who's responsible for each security function.

### 3. Risk Management Integration

Security governance and risk management go hand-in-hand. Your governance framework should include processes for:

- Identifying and assessing security risks
- Prioritizing risks based on business impact
- Implementing appropriate controls
- Monitoring and reporting on risk status

### 4. Compliance and Audit Programs

Regular audits ensure your security controls are working as intended and help you maintain compliance with regulations like GDPR, HIPAA, SOC 2, and others.

Use our [Security Audit Checklist Template](/templates/security-audit-checklist) to streamline your internal audit process.

## Benefits of Strong Security Governance

**Risk Reduction** - Proactive identification and mitigation of security threats before they become incidents.

**Regulatory Compliance** - Demonstrate compliance with industry regulations and avoid costly penalties.

**Stakeholder Confidence** - Build trust with customers, partners, and investors by showing commitment to security.

**Cost Efficiency** - Avoid redundant security efforts and optimize resource allocation.

**Incident Response** - Respond to security incidents faster and more effectively with clear procedures.

## Getting Started with Security Governance

1. **Assess Your Current State** - Understand what security measures you already have in place
2. **Define Your Security Strategy** - Align security goals with business objectives
3. **Create Core Policies** - Document essential security policies and procedures
4. **Assign Ownership** - Designate security champions across your organization
5. **Implement Monitoring** - Set up metrics to track security performance
6. **Review and Improve** - Regularly assess and update your governance framework

## Templates to Accelerate Your Governance Program

Building a security governance framework from scratch can be overwhelming. Our templates provide a head start:

- [Information Security Policy Template](/templates/information-security-policy)
- [Security Governance Charter](/templates/governance-charter)
- [Board Security Report Template](/templates/board-security-report)
- [Security Metrics Dashboard](/templates/security-metrics-dashboard)

## Conclusion

Security governance isn't a one-time project—it's an ongoing commitment to protecting your organization's most valuable assets. By establishing clear policies, defined responsibilities, and regular oversight, you create a security culture that grows stronger over time.

Ready to build your security governance framework? Start with our proven templates and customize them to fit your organization's unique needs.`,
			AuthorID:   1,
			CategoryID: 1,
			SEO:        "security governance, information security, security policies, governance framework, security management",
		},
		{
			Title:      "Risk Management Essentials: Identifying and Mitigating Security Threats",
			Content:    `# Risk Management Essentials: Identifying and Mitigating Security Threats

In today's digital landscape, every organization faces security risks—from data breaches and ransomware to insider threats and system vulnerabilities. The question isn't whether risks exist, but how effectively you identify, assess, and manage them.

## Why Risk Management Matters

Risk management is the systematic process of identifying potential security threats, evaluating their likelihood and impact, and implementing controls to reduce them to acceptable levels. It's about making informed decisions with limited resources.

Here's the reality: you can't eliminate all risks. But with a structured approach, you can:

- Focus resources on the most critical threats
- Make risk-informed business decisions
- Reduce the likelihood and impact of security incidents
- Demonstrate due diligence to stakeholders and regulators

## The Risk Management Process

### Step 1: Risk Identification

The first step is discovering what could go wrong. Common security risks include:

- **Technical Risks** - Vulnerabilities in software, hardware, or network infrastructure
- **Human Risks** - Phishing attacks, social engineering, insider threats
- **Process Risks** - Inadequate procedures, lack of access controls
- **External Risks** - Third-party vendors, supply chain vulnerabilities
- **Physical Risks** - Unauthorized facility access, environmental hazards

**Pro Tip:** Use our [Risk Identification Worksheet](/templates/risk-identification-worksheet) to conduct thorough risk discovery sessions with your team.

### Step 2: Risk Assessment

Once identified, each risk needs to be evaluated based on:

- **Likelihood** - How probable is this risk to occur?
- **Impact** - What would be the consequences if it happens?

Most organizations use a risk matrix to categorize risks as:
- **Critical** - Immediate action required
- **High** - Priority for mitigation
- **Medium** - Scheduled for remediation
- **Low** - Accept or monitor

Our [Risk Assessment Matrix Template](/templates/risk-assessment-matrix) provides a standardized framework for evaluating and prioritizing risks.

### Step 3: Risk Treatment

For each identified risk, you have four main options:

**Mitigate** - Implement controls to reduce likelihood or impact (most common)

**Accept** - Acknowledge the risk and choose not to take action (for low-priority risks)

**Transfer** - Shift the risk to another party (e.g., cyber insurance, outsourcing)

**Avoid** - Eliminate the risk by changing plans or processes

### Step 4: Risk Monitoring

Risk management isn't a one-time activity. Risks evolve as your organization and the threat landscape change. Establish ongoing monitoring through:

- Regular risk assessments (quarterly or annually)
- Continuous vulnerability scanning
- Security metrics and KPIs
- Incident tracking and analysis

Download our [Risk Register Template](/templates/risk-register) to maintain a living inventory of your organization's risks.

## Common Security Risks and How to Address Them

### Data Breaches

**Risk:** Unauthorized access to sensitive customer or business data

**Mitigation Strategies:**
- Implement encryption for data at rest and in transit
- Use multi-factor authentication (MFA)
- Apply principle of least privilege access
- Conduct regular security awareness training

### Ransomware Attacks

**Risk:** Malware that encrypts your data and demands payment

**Mitigation Strategies:**
- Maintain offline, encrypted backups
- Keep systems and software patched
- Deploy endpoint detection and response (EDR) tools
- Create and test incident response plans

Use our [Incident Response Plan Template](/templates/incident-response-plan) to prepare your team for security events.

### Third-Party Vendor Risks

**Risk:** Security weaknesses in vendor systems compromising your data

**Mitigation Strategies:**
- Conduct vendor security assessments
- Include security requirements in contracts
- Monitor vendor security posture
- Limit vendor access to only necessary systems

Our [Vendor Risk Assessment Template](/templates/vendor-risk-assessment) helps you evaluate third-party security risks systematically.

## Building a Risk-Aware Culture

Technology alone won't protect your organization—you need people who understand and care about security. Build a risk-aware culture by:

- **Training employees** on security best practices and threat recognition
- **Communicating openly** about risks and security incidents
- **Rewarding** employees who identify and report security concerns
- **Leading by example** with executive commitment to security

## Risk Management Tools and Templates

Streamline your risk management program with our professional templates:

- [Risk Management Policy](/templates/risk-management-policy)
- [Threat Assessment Template](/templates/threat-assessment-template)
- [Business Impact Analysis (BIA)](/templates/business-impact-analysis)
- [Risk Treatment Plan](/templates/risk-treatment-plan)
- [Third-Party Risk Questionnaire](/templates/third-party-risk-questionnaire)

## Measuring Risk Management Success

Track these key metrics to evaluate your risk management program:

- **Number of identified risks** being actively managed
- **Percentage of high/critical risks** mitigated within target timeframes
- **Mean time to remediate** vulnerabilities
- **Security incident trends** (increasing or decreasing)
- **Risk acceptance rate** by senior leadership

## Conclusion

Effective risk management is about making your organization more resilient. By systematically identifying threats, assessing their potential impact, and implementing appropriate controls, you transform security from a reactive cost center into a strategic business enabler.

Don't let security risks keep you up at night. Start building your risk management program today with our comprehensive templates and frameworks.`,
			AuthorID:   2,
			CategoryID: 1,
			SEO:        "risk management, security risk, threat assessment, risk assessment, cybersecurity risk, risk mitigation",
		},
		{
			Title:      "Compliance Made Simple: Navigating Security Regulations and Standards",
			Content:    `# Compliance Made Simple: Navigating Security Regulations and Standards

If you've ever felt overwhelmed by the alphabet soup of compliance requirements—GDPR, HIPAA, SOC 2, PCI DSS, ISO 27001—you're not alone. Compliance can feel like a moving target, but it doesn't have to be complicated.

Let's break down what compliance really means and how you can approach it strategically without drowning in documentation and bureaucracy.

## What is Compliance?

Compliance means adhering to laws, regulations, standards, and contractual obligations that apply to your organization. In the security context, it's about demonstrating that you're protecting sensitive data according to established requirements.

Think of compliance as the **minimum baseline** for security. It's not the finish line—it's the starting point for building a robust security program.

## Why Compliance Matters

**Legal Requirements** - Many industries have mandatory regulations. Non-compliance can result in fines, legal action, and business restrictions.

**Customer Trust** - Compliance demonstrates your commitment to protecting customer data, building confidence and competitive advantage.

**Business Opportunities** - Many enterprise customers won't work with vendors who lack certain certifications (like SOC 2 or ISO 27001).

**Risk Reduction** - Compliance frameworks guide you toward security best practices, reducing your overall risk exposure.

## Common Security Compliance Frameworks

### GDPR (General Data Protection Regulation)

**Applies to:** Organizations processing personal data of EU residents

**Key Requirements:**
- Obtain explicit consent for data processing
- Provide data access and deletion rights
- Report breaches within 72 hours
- Appoint Data Protection Officer (if applicable)
- Implement privacy by design

Get started with our [GDPR Compliance Checklist](/templates/gdpr-compliance-checklist) and [Privacy Policy Template](/templates/privacy-policy-gdpr).

### SOC 2 (Service Organization Control 2)

**Applies to:** Service providers storing customer data in the cloud

**Key Requirements:**
- Document security policies and procedures
- Implement access controls
- Monitor and log system activity
- Conduct regular risk assessments
- Provide audit evidence

Our [SOC 2 Readiness Assessment](/templates/soc2-readiness-assessment) and [SOC 2 Policy Bundle](/templates/soc2-policy-bundle) accelerate your certification journey.

### HIPAA (Health Insurance Portability and Accountability Act)

**Applies to:** Healthcare providers and their business associates

**Key Requirements:**
- Protect electronic Protected Health Information (ePHI)
- Conduct regular risk assessments
- Implement access controls and audit logs
- Provide breach notification
- Train workforce on privacy and security

Download our [HIPAA Compliance Toolkit](/templates/hipaa-compliance-toolkit) to ensure you're meeting all requirements.

### PCI DSS (Payment Card Industry Data Security Standard)

**Applies to:** Organizations that store, process, or transmit credit card data

**Key Requirements:**
- Maintain secure network and systems
- Protect cardholder data with encryption
- Implement strong access controls
- Monitor and test networks regularly
- Maintain information security policy

Use our [PCI DSS Compliance Guide](/templates/pci-dss-compliance-guide) to navigate the 12 requirements.

### ISO 27001

**Applies to:** Organizations wanting internationally recognized security certification

**Key Requirements:**
- Establish Information Security Management System (ISMS)
- Conduct risk assessments
- Implement security controls
- Document policies and procedures
- Perform internal audits

Our [ISO 27001 Implementation Kit](/templates/iso27001-implementation-kit) includes all required policies and procedures.

## Building Your Compliance Program

### 1. Identify Applicable Requirements

Start by determining which regulations and standards apply to your organization based on:

- **Industry** (healthcare, finance, e-commerce, etc.)
- **Geography** (where you operate and where your customers are located)
- **Data types** (personal data, payment information, health records)
- **Customer requirements** (contractual compliance obligations)

### 2. Perform a Gap Analysis

Assess your current security posture against compliance requirements. Identify:

- What you're already doing right
- Where gaps exist
- What documentation is missing
- Which controls need implementation

Our [Compliance Gap Analysis Template](/templates/compliance-gap-analysis) provides a structured approach to this assessment.

### 3. Create a Compliance Roadmap

Prioritize remediation efforts based on:

- **Risk level** (address critical gaps first)
- **Required timeline** (regulatory deadlines)
- **Resource availability** (staff, budget, technology)
- **Quick wins** (easy improvements with high impact)

### 4. Implement Controls and Documentation

Compliance requires both technical controls and documentation:

**Technical Controls:**
- Access management systems
- Encryption solutions
- Monitoring and logging tools
- Backup and recovery systems

**Documentation:**
- Policies and procedures
- Risk assessments
- Training records
- Audit evidence

Download our [Compliance Documentation Bundle](/templates/compliance-documentation-bundle) to streamline documentation requirements.

### 5. Train Your Team

Your employees are the front line of compliance. Ensure everyone understands:

- Their role in maintaining compliance
- How to handle sensitive data
- How to recognize and report security incidents
- Why compliance matters

### 6. Monitor and Maintain

Compliance isn't a one-time achievement—it requires ongoing effort:

- Conduct regular internal audits
- Update policies as regulations change
- Review access controls quarterly
- Track and respond to security incidents
- Prepare for external audits

Use our [Compliance Monitoring Dashboard](/templates/compliance-monitoring-dashboard) to track your compliance posture.

## Multi-Framework Compliance

Many organizations must comply with multiple frameworks simultaneously. The good news? There's significant overlap:

- Access controls appear in nearly every framework
- Risk assessments are universally required
- Incident response is a common requirement
- Documentation and training are standard expectations

**Strategy:** Implement comprehensive security controls that satisfy multiple frameworks rather than treating each as separate. Our [Unified Compliance Matrix](/templates/unified-compliance-matrix) maps controls across major frameworks.

## Common Compliance Pitfalls to Avoid

**Checkbox Mentality** - Don't just check boxes to pass audits. Build security that actually protects your organization.

**Documentation Overload** - Focus on useful documentation that guides action, not volumes of paperwork nobody reads.

**Point-in-Time Compliance** - Compliance is continuous, not just something you do before an audit.

**Ignoring Third Parties** - Your vendors must also be compliant. Their security failures become your compliance failures.

**Lack of Executive Buy-In** - Compliance requires resources. Ensure leadership understands the business value.

## Compliance Templates and Resources

Accelerate your compliance journey with our expert-designed templates:

- [Compliance Program Charter](/templates/compliance-program-charter)
- [Policy Template Library](/templates/policy-template-library)
- [Audit Response Template](/templates/audit-response-template)
- [Control Testing Workbook](/templates/control-testing-workbook)
- [Compliance Training Materials](/templates/compliance-training-materials)
- [Evidence Collection Guide](/templates/evidence-collection-guide)

## Working with Auditors

When audit time comes, be prepared:

- **Organize evidence** in advance (screenshots, logs, training records)
- **Be honest** about gaps—auditors appreciate transparency
- **Ask questions** to understand requirements better
- **Take notes** on findings and recommendations
- **Follow up** on remediation commitments

Our [Audit Preparation Checklist](/templates/audit-preparation-checklist) ensures you're ready for external assessments.

## Compliance as Competitive Advantage

Rather than viewing compliance as a burden, consider it an investment:

- **Win more customers** who require security certifications
- **Command premium pricing** for demonstrable security
- **Reduce breach risk** through proven security controls
- **Streamline operations** with documented processes
- **Attract better talent** who value security-conscious organizations

## Conclusion

Compliance doesn't have to be overwhelming. With a systematic approach, the right templates, and ongoing commitment, you can achieve and maintain compliance while building genuine security.

Remember: compliance is the floor, not the ceiling. Use regulatory requirements as a foundation, then build additional security measures that make sense for your specific risk profile.

Ready to tackle your compliance requirements? Browse our comprehensive library of compliance templates and start building your program today.`,
			AuthorID:   1,
			CategoryID: 1,
			SEO:        "compliance, security compliance, GDPR, SOC 2, HIPAA, PCI DSS, ISO 27001, regulatory compliance",
		},
	}

	for _, post := range posts {
		if err := s.db.Create(&post).Error; err != nil {
			if !errors.Is(err, gorm.ErrDuplicatedKey) {
				return err
			}
		}
	}

	return nil
}