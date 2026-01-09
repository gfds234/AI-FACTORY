import React from 'react'

export default function Hero() {
  return (
    <section className="hero">
      <div className="hero-content">
        <h1 className="hero-title">
          Streamline Your Projects with <span className="gradient-text">TaskFlow Pro</span>
        </h1>
        <p className="hero-subtitle">
          The modern project management tool that helps teams collaborate, track progress, and deliver projects on time.
        </p>
        <div className="hero-cta">
          <button className="btn btn-primary">Start Free Trial</button>
          <button className="btn btn-secondary">Watch Demo</button>
        </div>
        <div className="hero-stats">
          <div className="stat">
            <div className="stat-number">50K+</div>
            <div className="stat-label">Active Users</div>
          </div>
          <div className="stat">
            <div className="stat-number">99.9%</div>
            <div className="stat-label">Uptime</div>
          </div>
          <div className="stat">
            <div className="stat-number">4.9/5</div>
            <div className="stat-label">User Rating</div>
          </div>
        </div>
      </div>
    </section>
  )
}
