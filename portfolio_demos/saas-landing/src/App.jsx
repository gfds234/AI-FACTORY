import { useState } from 'react'
import './App.css'
import Hero from './components/Hero'
import Features from './components/Features'
import Pricing from './components/Pricing'
import Contact from './components/Contact'

function App() {
  return (
    <div className="App">
      <Hero />
      <Features />
      <Pricing />
      <Contact />
      <footer className="footer">
        <p>&copy; 2024 TaskFlow Pro. All rights reserved.</p>
      </footer>
    </div>
  )
}

export default App
