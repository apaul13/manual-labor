import './App.css'
import { FaCar, FaHome, FaSearch } from "react-icons/fa";

function App() {
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
            <div className="search-option">
              <label htmlFor="vin-input" className="label-category-select">Add VIN</label>

                <input type="text" id="vin-input" className="input vin" placeholder="Enter VIN" />
            </div>
            <div className="search-option">
              <label htmlFor="year-input" className="label-category-select">Year</label>
              <div className="horizontal-divider">
                <select id="year-input" className="input">
                  <option value="">Select Year</option>
                  <option value="2024">2024</option>
                  <option value="2023">2023</option>
                  <option value="2022">2022</option>
                </select>
              </div>
              <label htmlFor="make-input" className="label-category-select">Make</label>
              <div className="horizontal-divider">
                <select id="make-input" className="input">
                  <option value="">Select Make</option>
                  <option value="honda">Honda</option>
                  <option value="toyota">Toyota</option>
                  <option value="ford">Ford</option>
                </select>
              </div>
              <label htmlFor="model-input" className="label-category-select">Model</label>
              <div className="horizontal-divider">
                <select id="model-input" className="input">
                  <option value="">Select Model</option>
                  <option value="model1">Model 1</option>
                  <option value="model2">Model 2</option>
                  <option value="model3">Model 3</option>
                </select>
              </div>
              <label htmlFor="trim-input" className="label-category-select">Trim</label>
              <div className="horizontal-divider">
                <select id="trim-input" className="input">
                  <option value="">Select Trim</option>
                  <option value="trim1">Trim 1</option>
                  <option value="trim2">Trim 2</option>
                  <option value="trim3">Trim 3</option>
                </select>
              </div>
            </div>
          </div>
          <div className="search-box">
            <button className="search-button"><FaSearch /> Search</button>
          </div>
        </section>
      </main>
    </>
  )
}

export default App
