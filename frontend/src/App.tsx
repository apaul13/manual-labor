import './App.css'
import { useState } from 'react'
import { FaCar, FaHome, FaSearch } from "react-icons/fa";
import { config } from './constants/config';

function App() {

  const [mode, setMode] = useState<'vin'|'details'>('vin')
  const [vin, setVin] = useState('')
  const [year, setYear] = useState('')
  const [make, setMake] = useState('')
  const [model, setModel] = useState('')
  const [trim, setTrim] = useState('')

  const canSearch = () => {
    if (mode === 'vin') return vin.trim().length > 0
    return Boolean(year || make || model || trim)
  }

  const handleSearch = () => {
    if (!canSearch()) return

    let url = `${config.API_BASE_URL}/cars`
    if (mode === 'vin') {
      url += `?vin=${encodeURIComponent(vin.trim())}`
    } else {
      const params = new URLSearchParams()
      if (year) params.append('year', year)
      if (make) params.append('make', make)
      if (model) params.append('model', model)
      if (trim) params.append('trim', trim)
      const qs = params.toString()
      if (qs) url += `?${qs}`
    }

    fetch(url)
      .then(response => response.json())
      .then(data => {
        console.log('Cars found:', data);
      })
      .catch(err => console.error(err))
  }

  return (
    <>
      <header>
        <nav className="navbar">
          <a href="/" className="nav-logo fugaz-one-regular">
            manny
          </a>
          <div className="nav-links">
            <a href="/auto" className="nav-link">
              <FaCar className="nav-icon" />
              <span>Auto</span>
            </a>
            <a href="/home" className="nav-link">
              <FaHome className="nav-icon" />
              <span>Home</span>
            </a>
          </div>
        </nav>
      </header>
      <main id="main-content">
        <section className="vehicle-identification">
          <div className="search-options">
            <div style={{gridColumn: '1 / -1'}}>
              <div className="mode-toggle" role="tablist" aria-label="Search mode">
                <button
                  className={`mode-button ${mode === 'vin' ? 'active' : ''}`}
                  onClick={() => setMode('vin')}
                >
                  VIN
                </button>
                <button
                  className={`mode-button ${mode === 'details' ? 'active' : ''}`}
                  onClick={() => setMode('details')}
                >
                  Details
                </button>
              </div>
            </div>

            {mode === 'vin' && (
              <div className={`field ${mode === 'vin' ? 'vin' : ''}`}>
                <label htmlFor="vin-input" className="label-category-select">VIN</label>
                <input
                  type="text"
                  id="vin-input"
                  className="input vin"
                  placeholder="Enter VIN (17 characters)"
                  value={vin}
                  onChange={e => setVin(e.target.value)}
                />
              </div>
            )}

            {mode === 'details' && (
              <>
                <div className="field">
                  <label htmlFor="year-input" className="label-category-select">Year</label>
                  <select id="year-input" className="input" value={year} onChange={e => setYear(e.target.value)}>
                    <option value="">Year</option>
                    <option value="2024">2024</option>
                    <option value="2023">2023</option>
                    <option value="2022">2022</option>
                  </select>
                </div>

                <div className="field">
                  <label htmlFor="make-input" className="label-category-select">Make</label>
                  <select id="make-input" className="input" value={make} disabled={!year} onChange={e => setMake(e.target.value)}>
                    <option value="">Make</option>
                    <option value="honda">Honda</option>
                    <option value="toyota">Toyota</option>
                    <option value="ford">Ford</option>
                  </select>
                </div>

                <div className="field">
                  <label htmlFor="model-input" className="label-category-select">Model</label>
                  <select id="model-input" className="input" value={model} disabled={!make} onChange={e => setModel(e.target.value)}>
                    <option value="">Model</option>
                    <option value="model1">Model 1</option>
                    <option value="model2">Model 2</option>
                    <option value="model3">Model 3</option>
                  </select>
                </div>

                <div className="field">
                  <label htmlFor="trim-input" className="label-category-select">Trim</label>
                  <select id="trim-input" className="input" value={trim} disabled={!model} onChange={e => setTrim(e.target.value)}>
                    <option value="">Trim</option>
                    <option value="trim1">Trim 1</option>
                    <option value="trim2">Trim 2</option>
                    <option value="trim3">Trim 3</option>
                  </select>
                </div>
              </>
            )}

          </div>

          <div className="search-box">
            <button
              className="search-button"
              onClick={handleSearch}
              aria-label="Search vehicles"
              disabled={!canSearch()}
              style={{opacity: canSearch() ? 1 : 0.6, cursor: canSearch() ? 'pointer' : 'not-allowed'}}
            >
              <FaSearch />
              <span style={{marginLeft:6}}>Search</span>
            </button>
          </div>
        </section>
      </main>
    </>
  )
}

export default App;