import React from 'react'

export default function Features() {
  const features = [
    {
      icon: '📋',
      title: 'Task Management',
      description: 'Create, assign, and track tasks with our intuitive kanban board interface. Stay organized and never miss a deadline.'
    },
    {
      icon: '👥',
      title: 'Team Collaboration',
      description: 'Work together seamlessly with real-time updates, comments, and file sharing. Keep everyone on the same page.'
    },
    {
      icon: '📊',
      title: 'Analytics Dashboard',
      description: 'Get insights into team productivity, project progress, and resource allocation with beautiful, actionable dashboards.'
    }
  ]

  return (
    <section className="features" id="features">
      <div className="container">
        <h2 className="section-title">Powerful Features for Modern Teams</h2>
        <p className="section-subtitle">Everything you need to manage projects efficiently</p>

        <div className="features-grid">
          {features.map((feature, index) => (
            <div key={index} className="feature-card">
              <div className="feature-icon">{feature.icon}</div>
              <h3 className="feature-title">{feature.title}</h3>
              <p className="feature-description">{feature.description}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}
