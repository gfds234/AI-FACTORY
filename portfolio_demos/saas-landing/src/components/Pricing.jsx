import React from 'react'

export default function Pricing() {
  const plans = [
    {
      name: 'Free',
      price: '$0',
      period: '/month',
      features: [
        'Up to 5 team members',
        'Basic task management',
        '10 projects',
        'Mobile app access',
        'Community support'
      ],
      highlighted: false
    },
    {
      name: 'Pro',
      price: '$29',
      period: '/month',
      features: [
        'Up to 25 team members',
        'Advanced task management',
        'Unlimited projects',
        'Priority support',
        'Custom workflows',
        'Analytics & reporting'
      ],
      highlighted: true
    },
    {
      name: 'Enterprise',
      price: '$99',
      period: '/month',
      features: [
        'Unlimited team members',
        'Everything in Pro',
        'Dedicated account manager',
        'Advanced security',
        'Custom integrations',
        'SLA guarantee'
      ],
      highlighted: false
    }
  ]

  return (
    <section className="pricing" id="pricing">
      <div className="container">
        <h2 className="section-title">Simple, Transparent Pricing</h2>
        <p className="section-subtitle">Choose the plan that fits your team</p>

        <div className="pricing-grid">
          {plans.map((plan, index) => (
            <div key={index} className={`pricing-card ${plan.highlighted ? 'highlighted' : ''}`}>
              {plan.highlighted && <div className="badge">Most Popular</div>}
              <h3 className="plan-name">{plan.name}</h3>
              <div className="plan-price">
                <span className="price">{plan.price}</span>
                <span className="period">{plan.period}</span>
              </div>
              <ul className="plan-features">
                {plan.features.map((feature, idx) => (
                  <li key={idx}>
                    <span className="checkmark">✓</span> {feature}
                  </li>
                ))}
              </ul>
              <button className={`btn ${plan.highlighted ? 'btn-primary' : 'btn-secondary'}`}>
                Get Started
              </button>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}
