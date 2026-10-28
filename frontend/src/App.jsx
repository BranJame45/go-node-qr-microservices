import { useState } from 'react'

const API_URL = import.meta.env.VITE_GO_API_URL || 'http://localhost:8080'

// Formatea números con hasta 4 decimales, sin ceros sobrantes
function fmt(val) {
  if (typeof val !== 'number' || !isFinite(val)) return String(val)
  return Number(val.toFixed(4)).toString()
}

function MatrixTable({ matrix }) {
  if (!Array.isArray(matrix)) return null
  return (
    <div className="matrix-wrap">
      <table className="matrix">
        <tbody>
          {matrix.map((row, i) => (
            <tr key={i}>
              {Array.isArray(row) ? row.map((val, j) => (
                <td key={j}>{typeof val === 'number' ? fmt(val) : val}</td>
              )) : null}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

function StatsGrid({ stats }) {
  if (!stats) return null
  return (
    <div className="stats-grid">
      <div className="stat"><p className="label">Máximo</p><div className="value">{fmt(stats.max)}</div></div>
      <div className="stat"><p className="label">Mínimo</p><div className="value">{fmt(stats.min)}</div></div>
      <div className="stat"><p className="label">Promedio</p><div className="value">{fmt(stats.average)}</div></div>
      <div className="stat"><p className="label">Suma</p><div className="value">{fmt(stats.sum)}</div></div>
      <div className="stat">
        <p className="label">Diagonal</p>
        <div className="value">
          <span className={`badge ${stats.isDiagonal ? 'yes' : 'no'}`}>
            {stats.isDiagonal ? 'Sí' : 'No'}
          </span>
        </div>
      </div>
    </div>
  )
}

function SummaryRow({ name, stats }) {
  if (!stats) return null
  return (
    <tr>
      <td>{name}</td>
      <td>{fmt(stats.max)}</td>
      <td>{fmt(stats.min)}</td>
      <td>{fmt(stats.average)}</td>
      <td>{fmt(stats.sum)}</td>
      <td><span className={`badge ${stats.isDiagonal ? 'yes' : 'no'}`}>{stats.isDiagonal ? 'Sí' : 'No'}</span></td>
    </tr>
  )
}

function App() {
  const [token, setToken] = useState('')
  const [username, setUsername] = useState('admin')
  const [password, setPassword] = useState('admin123')
  const [matrixInput, setMatrixInput] = useState('[[12,-51,4],[6,167,-68],[-4,24,-41]]')
  const [result, setResult] = useState(null)
  const [error, setError] = useState('')
  const [loadingLogin, setLoadingLogin] = useState(false)
  const [loadingCalc, setLoadingCalc] = useState(false)

  const handleLogin = async (e) => {
    e.preventDefault()
    setError('')
    setLoadingLogin(true)
    try {
      const res = await fetch(`${API_URL}/api/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password })
      })
      const data = await res.json()
      if (res.ok) setToken(data.token)
      else setError(data.error || 'Error de autenticación')
    } catch {
      setError('No se pudo conectar con la API')
    } finally {
      setLoadingLogin(false)
    }
  }

  const handleLogout = () => {
    setToken('')
    setResult(null)
    setError('')
  }

  const handleCalculate = async (e) => {
    e.preventDefault()
    setError('')
    setResult(null)
    let matrix
    try {
      matrix = JSON.parse(matrixInput)
    } catch {
      setError('La matriz no tiene un formato JSON válido. Ejemplo: [[1,2],[3,4]]')
      return
    }
    setLoadingCalc(true)
    try {
      const res = await fetch(`${API_URL}/api/qr`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({ matrix })
      })
      const data = await res.json()
      if (res.ok) setResult(data)
      else setError(data.error || 'Error al calcular la factorización')
    } catch {
      setError('No se pudo conectar con la API')
    } finally {
      setLoadingCalc(false)
    }
  }

  return (
    <>
      <header className="app-header">
        <div className="app-header-inner">
          <div className="brand">
            <h1>Factorización QR</h1>
            <p>Cálculo de matrices Q y R con análisis estadístico</p>
          </div>
          {token && (
            <div className="session">
              <span><span className="dot" />Sesión activa</span>
              <button className="ghost" onClick={handleLogout}>Cerrar sesión</button>
            </div>
          )}
        </div>
      </header>

      {!token ? (
        <div className="container narrow">
          <div className="card">
            <h2>Iniciar sesión</h2>
            <p className="hint">Ingresa tus credenciales para acceder al servicio.</p>
            <form onSubmit={handleLogin}>
              <label htmlFor="user">Usuario</label>
              <input id="user" type="text" value={username} onChange={(e) => setUsername(e.target.value)} />
              <label htmlFor="pass">Contraseña</label>
              <input id="pass" type="password" value={password} onChange={(e) => setPassword(e.target.value)} />
              <button type="submit" disabled={loadingLogin}>
                {loadingLogin ? 'Ingresando…' : 'Iniciar sesión'}
              </button>
            </form>
            {error && <div className="alert" style={{ marginTop: 18, marginBottom: 0 }}>{error}</div>}
          </div>
        </div>
      ) : (
        <div className="container">
          <div className="card">
            <h2>Matriz de entrada</h2>
            <p className="hint">Introduce una matriz rectangular (m ≥ n) en formato JSON.</p>
            <form onSubmit={handleCalculate}>
              <textarea
                rows="5"
                value={matrixInput}
                onChange={(e) => setMatrixInput(e.target.value)}
                placeholder="[[1,2],[3,4]]"
              />
              <button type="submit" disabled={loadingCalc}>
                {loadingCalc ? 'Calculando…' : 'Calcular factorización'}
              </button>
            </form>
          </div>

          {error && <div className="alert">{error}</div>}

          {result && (
            <>
              <p className="section-title">Resultado de la factorización</p>
              <div className="card">
                <h2>Matriz Q <span style={{ fontWeight: 400, color: 'var(--text-muted)', fontSize: '0.82rem' }}>· ortogonal</span></h2>
                <MatrixTable matrix={result.Q} />
              </div>
              <div className="card">
                <h2>Matriz R <span style={{ fontWeight: 400, color: 'var(--text-muted)', fontSize: '0.82rem' }}>· triangular superior</span></h2>
                <MatrixTable matrix={result.R} />
              </div>

              {result.statistics && (
                <>
                  <p className="section-title">Análisis estadístico</p>
                  <div className="card">
                    <h2>Global</h2>
                    <p className="hint">Sobre el conjunto de valores de Q y R.</p>
                    <StatsGrid stats={result.statistics.global} />

                    {result.statistics.perMatrix && (
                      <div className="per-matrix">
                        <h4>Desglose por matriz</h4>
                        <table className="summary">
                          <thead>
                            <tr>
                              <th>Matriz</th><th>Máx</th><th>Mín</th><th>Promedio</th><th>Suma</th><th>Diagonal</th>
                            </tr>
                          </thead>
                          <tbody>
                            <SummaryRow name="Q" stats={result.statistics.perMatrix.Q} />
                            <SummaryRow name="R" stats={result.statistics.perMatrix.R} />
                          </tbody>
                        </table>
                      </div>
                    )}
                  </div>
                </>
              )}
            </>
          )}
        </div>
      )}
    </>
  )
}

export default App
