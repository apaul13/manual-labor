import './App.css'
import { FaCar, FaHome } from "react-icons/fa";
import CarSearch from './components/CarSearch/CarSearch';

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
        <CarSearch />
      </main>
    </>
  )
}

export default App;