import React, { useState } from 'react'

export default function Contact() {
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    message: ''
  })
  const [submitted, setSubmitted] = useState(false)

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    })
  }

  const handleSubmit = (e) => {
    e.preventDefault()
    // In production, this would send to an API
    console.log('Form submitted:', formData)
    setSubmitted(true)
    setTimeout(() => {
      setSubmitted(false)
      setFormData({ name: '', email: '', message: '' })
    }, 3000)
  }

  return (
    <section className="contact" id="contact">
      <div className="container">
        <h2 className="section-title">Get in Touch</h2>
        <p className="section-subtitle">Have questions? We'd love to hear from you.</p>

        <div className="contact-content">
          <form className="contact-form" onSubmit={handleSubmit}>
            <div className="form-group">
              <label htmlFor="name">Name</label>
              <input
                type="text"
                id="name"
                name="name"
                value={formData.name}
                onChange={handleChange}
                placeholder="Your name"
                required
              />
            </div>

            <div className="form-group">
              <label htmlFor="email">Email</label>
              <input
                type="email"
                id="email"
                name="email"
                value={formData.email}
                onChange={handleChange}
                placeholder="your.email@example.com"
                required
              />
            </div>

            <div className="form-group">
              <label htmlFor="message">Message</label>
              <textarea
                id="message"
                name="message"
                value={formData.message}
                onChange={handleChange}
                placeholder="Tell us about your project..."
                rows="5"
                required
              />
            </div>

            <button type="submit" className="btn btn-primary">
              {submitted ? 'Message Sent! ✓' : 'Send Message'}
            </button>

            {submitted && (
              <p className="success-message">Thank you! We'll get back to you soon.</p>
            )}
          </form>

          <div className="contact-info">
            <div className="info-item">
              <span className="info-icon">📧</span>
              <div>
                <h4>Email</h4>
                <p>hello@taskflowpro.com</p>
              </div>
            </div>
            <div className="info-item">
              <span className="info-icon">💬</span>
              <div>
                <h4>Live Chat</h4>
                <p>Available 24/7</p>
              </div>
            </div>
            <div className="info-item">
              <span className="info-icon">📍</span>
              <div>
                <h4>Office</h4>
                <p>San Francisco, CA</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  )
}
