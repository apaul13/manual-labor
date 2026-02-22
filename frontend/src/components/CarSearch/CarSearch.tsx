
import { useState } from 'react'
import { FaSearch } from "react-icons/fa";
import { config } from '../../constants/config';
import './CarSearch.css'

export default function CarSearch() {

    const [mode, setMode] = useState<'vin' | 'details'>('vin')
    const [vin, setVin] = useState('')
    const [year, setYear] = useState('')
    const [make, setMake] = useState('')
    const [model, setModel] = useState('')
    const [trim, setTrim] = useState('')

    const canSearch = () => {
        if (mode === 'vin')
            return vin.trim().length > 0
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
        <section className="car-search">
            <form className="vehicle-identification" onSubmit={(e) => { e.preventDefault(); handleSearch(); }}>
                <div className="search-options">
                    <div style={{ gridColumn: '1 / -1' }}>
                        <div className="mode-toggle" role="tablist" aria-label="Search mode">
                            <button
                                className={`mode-button ${mode === 'vin' ? 'active' : ''}`}
                                type="button"
                                onClick={() => setMode('vin')}
                            >
                                VIN
                            </button>
                            <button
                                className={`mode-button ${mode === 'details' ? 'active' : ''}`}
                                type="button"
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
                                pattern="[A-HJ-NPR-Z0-9]{17}"
                                value={vin}
                                autoComplete='false'
                                onChange={e => {
                                    setVin(e.target.value);
                                    e.target.setCustomValidity("");
                                }}
                                onInvalid={e => {
                                    let target = e.target as HTMLInputElement;
                                    target.setCustomValidity("VIN must be exactly 17 characters and contain only valid characters.");
                                }}

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
                                    <option value="2021">2021</option>
                                    <option value="2020">2020</option>
                                </select>
                            </div>

                            <div className="field">
                                <label htmlFor="make-input" className="label-category-select">Make</label>
                                <select id="make-input" className="input" value={make} disabled={!year} onChange={e => setMake(e.target.value)}>
                                    <option value="">Make</option>
                                    <option value="HONDA">Honda</option>
                                    <option value="TOYOTA">Toyota</option>
                                    <option value="FORD">Ford</option>
                                </select>
                            </div>

                            <div className="field">
                                <label htmlFor="model-input" className="label-category-select">Model</label>
                                <select id="model-input" className="input" value={model} disabled={!make} onChange={e => setModel(e.target.value)}>
                                    <option value="">Model</option>
                                    <option value="PRIUS">Prius</option>
                                    <option value="CAMRY">Camry</option>
                                    <option value="COROLLA">Corolla</option>
                                </select>
                            </div>

                            <div className="field">
                                <label htmlFor="trim-input" className="label-category-select">Trim</label>
                                <select id="trim-input" className="input" value={trim} disabled={!model} onChange={e => setTrim(e.target.value)}>
                                    <option value="">Trim</option>
                                    <option value="SE">SE</option>
                                    <option value="LE">LE</option>
                                    <option value="XLE">XLE</option>
                                </select>
                            </div>
                        </>
                    )}

                </div>

                <div className="search-box">
                    <button
                        className="search-button"
                        aria-label="Search vehicles"
                        disabled={!canSearch()}
                        style={{ opacity: canSearch() ? 1 : 0.6, cursor: canSearch() ? 'pointer' : 'not-allowed' }}
                    >
                        <FaSearch />
                        <span style={{ marginLeft: 6 }}>Search</span>
                    </button>
                </div>
            </form>
        </section>
    )
}