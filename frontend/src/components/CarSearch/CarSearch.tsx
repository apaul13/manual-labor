
import { useEffect, useState } from 'react'
import { FaSearch } from "react-icons/fa";
import { config } from '../../constants/config';
import './CarSearch.css'
import type { Make, Model, Trim, MakesResponse, YearsResponse, Year } from '../../types/models';

export default function CarSearch() {

    const [mode, setMode] = useState<'vin' | 'details'>('vin')
    const [vin, setVin] = useState('')
    const [year, setYear] = useState('')
    const [make, setMake] = useState('')
    const [model, setModel] = useState('')
    const [trim, setTrim] = useState('')

    const [makes, setMakes] = useState<Make[]>([])
    const [years, setYears] = useState<string[]>([])
    const [models, setModels] = useState<Model[]>([])
    const [trims, setTrims] = useState<Trim[]>([])

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

    const handleSelectYear = (e: React.ChangeEvent<HTMLSelectElement>) => {
        setYear(e.target.value);
        setMake('');
        setModel('');
        setTrim('');

        fetch(`${config.API_BASE_URL}/makes?year=${e.target.value}&offset=0&limit=1000`)
            .then(response => response.json() as Promise<MakesResponse>)
            .then(data => {
                if (!data ||!Array.isArray(data?.Makes)) {
                    console.error('Expected an array of makes, but got:', data);
                    return;
                }
                console.log('Makes for year', e.target.value, ':', data);
                setMakes(data.Makes);
            })
            .catch(err => console.error(err))
    }

    const handleSelectMake = (e: React.ChangeEvent<HTMLSelectElement>) => {
        setMake(e.target.value);
        setModel('');
        setTrim('');

        fetch(`${config.API_BASE_URL}/models?make=${e.target.value}&offset=0&limit=1000`)
            .then(response => response.json())
            .then(data => {
                if (!Array.isArray(data)) {
                    console.error('Expected an array of models, but got:', data);
                    return;
                }
                console.log('Models for make', e.target.value, ':', data);
                setModels(data);
            })
            .catch(err => console.error(err)
            )
    }

    const handleSelectModel = (e: React.ChangeEvent<HTMLSelectElement>) => {
        setModel(e.target.value);
        setTrim('');

        fetch(`${config.API_BASE_URL}/trims?make=${make}&model=${e.target.value}&offset=0&limit=1000`)
            .then(response => response.json())
            .then(data => {
                if (!Array.isArray(data)) {
                    console.error('Expected an array of trims, but got:', data);
                    return;
                }
                console.log('Trims for make', make, 'model', e.target.value, ':', data);
                setTrims(data);
            })
            .catch(err => console.error(err)
            )
    }

    useEffect(() => {
        fetch(`${config.API_BASE_URL}/years?offset=0&limit=1000`)
            .then(response => response.json() as Promise<YearsResponse>)
            .then(data => {
                if (!data || !Array.isArray(data?.Years)) {
                    console.error('Expected an array of makes, but got:', data);
                    return;
                }
                setYears(Array.from(new Set(data.Years.map((y: Year) => y.Year))).sort((a, b) => parseInt(b) - parseInt(a)));
            })
            .catch(err => console.error(err)
            )
    }, [])

    console.log('Current state:', { mode, vin, year, make, model, trim, makes, years });

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
                                <select id="year-input" className="input" value={year} onChange={handleSelectYear}>
                                    <option value="">Year</option>
                                    {years && years.map(y => (
                                        <option key={y} value={y}>{y}</option>
                                    ))}
                                </select>
                            </div>

                            <div className="field">
                                <label htmlFor="make-input" className="label-category-select">Make</label>
                                <select id="make-input" className="input" value={make} disabled={!year} onChange={handleSelectMake}>
                                    <option value="">Make</option>
                                    {makes && makes.map(m => (
                                        <option key={m.ID} value={m.ID}>{m.Name}</option>
                                    ))}
                                </select>
                            </div>

                            <div className="field">
                                <label htmlFor="model-input" className="label-category-select">Model</label>
                                <select id="model-input" className="input" value={model} disabled={!make} onChange={handleSelectModel}>
                                    <option value="">Model</option>
                                    {models && models.map(m => (
                                        <option key={m.ID} value={m.ID}>{m.Name}</option>
                                    ))}
                                </select>
                            </div>

                            <div className="field">
                                <label htmlFor="trim-input" className="label-category-select">Trim</label>
                                <select id="trim-input" className="input" value={trim} disabled={!model} onChange={e => setTrim(e.target.value)}>
                                    <option value="">Trim</option>
                                    {trims && trims.map(t => (
                                        <option key={t.ID} value={t.ID}>{t.Name}</option>
                                    ))}
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