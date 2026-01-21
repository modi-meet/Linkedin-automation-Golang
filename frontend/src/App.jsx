import { useState, useEffect, useRef } from 'react'

function App() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [keyword, setKeyword] = useState('')
  const [limit, setLimit] = useState(5)
  const [message, setMessage] = useState('')
  const [logs, setLogs] = useState([])
  const [isRunning, setIsRunning] = useState(false)
  const [showCreds, setShowCreds] = useState(false)
  const logsEndRef = useRef(null)

  useEffect(() => {
    if (isRunning) {
      const eventSource = new EventSource('http://localhost:8080/api/events')

      eventSource.onmessage = (event) => {
        setLogs((prev) => [...prev, { text: event.data, time: new Date().toLocaleTimeString() }])
      }

      eventSource.onerror = () => {
        eventSource.close()
        setIsRunning(false)
      }

      return () => {
        eventSource.close()
      }
    }
  }, [isRunning])

  useEffect(() => {
    logsEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [logs])

  const handleStart = async () => {
    setLogs([])
    setIsRunning(true)

    try {
      const res = await fetch('http://localhost:8080/api/start', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          Email: email,
          Password: password,
          Keyword: keyword,
          Limit: parseInt(limit),
          ConnectMessage: message,
          Headless: false
        })
      })

      if (!res.ok) {
        throw new Error('Failed to start automation')
      }
    } catch (err) {
      setLogs(prev => [...prev, { text: `Error: ${err.message}`, time: new Date().toLocaleTimeString(), isError: true }])
      setIsRunning(false)
    }
  }

  const getLogStyle = (log) => {
    if (log.isError) return 'text-red-600'
    if (log.text.includes('successful') || log.text.includes('complete')) return 'text-emerald-600'
    if (log.text.includes('Skipping') || log.text.includes('Warning')) return 'text-amber-600'
    return 'text-slate-700'
  }

  return (
    <div className="min-h-screen bg-slate-50 text-slate-800 p-6 md:p-10">
      <div className="max-w-5xl mx-auto space-y-6">

        {/* Header */}
        <div className="flex items-center justify-between pb-4 border-b border-slate-200">
          <div className="flex items-center gap-3">
            <div className="w-9 h-9 rounded-lg bg-gradient-to-br from-blue-500 to-blue-600 flex items-center justify-center shadow-sm">
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className="text-white">
                <path d="M16 8a6 6 0 0 1 6 6v7h-4v-7a2 2 0 0 0-2-2 2 2 0 0 0-2 2v7h-4v-7a6 6 0 0 1 6-6z"></path>
                <rect x="2" y="9" width="4" height="12"></rect>
                <circle cx="4" cy="4" r="2"></circle>
              </svg>
            </div>
            <div>
              <h1 className="text-lg font-semibold text-slate-900 tracking-tight">LinkedIn Automation</h1>
              <p className="text-xs text-slate-500">Stealth connection workflow</p>
            </div>
          </div>

          <div className="flex items-center gap-2 px-3 py-1.5 bg-white rounded-full border border-slate-200 shadow-sm">
            <div className={`w-2 h-2 rounded-full ${isRunning ? 'bg-emerald-500 animate-pulse' : 'bg-slate-300'}`}></div>
            <span className="text-xs text-slate-500 font-medium">{isRunning ? 'Running' : 'Idle'}</span>
          </div>
        </div>

        {/* Controls */}
        <div className="grid grid-cols-1 md:grid-cols-12 gap-4">

          {/* Credentials Section */}
          <div className="md:col-span-3">
            <button
              onClick={() => setShowCreds(!showCreds)}
              className="w-full flex items-center justify-between gap-2 px-4 py-3 bg-white border border-slate-200 rounded-xl hover:border-slate-300 hover:shadow-sm transition-all text-sm"
            >
              <div className="flex items-center gap-2">
                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className="text-slate-400">
                  <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
                  <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
                </svg>
                <span className="text-slate-600">Credentials</span>
              </div>
              <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className={`text-slate-400 transition-transform ${showCreds ? 'rotate-180' : ''}`}>
                <polyline points="6 9 12 15 18 9"></polyline>
              </svg>
            </button>

            {showCreds && (
              <div className="mt-2 p-4 bg-white border border-slate-200 rounded-xl shadow-sm space-y-3">
                <div>
                  <label className="block text-xs text-slate-500 mb-1.5 font-medium">Email</label>
                  <input
                    type="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    className="w-full bg-slate-50 border border-slate-200 rounded-lg px-3 py-2 text-sm focus:border-blue-400 focus:ring-2 focus:ring-blue-100 outline-none transition-all placeholder-slate-400"
                    placeholder="user@example.com"
                  />
                </div>
                <div>
                  <label className="block text-xs text-slate-500 mb-1.5 font-medium">Password</label>
                  <input
                    type="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    className="w-full bg-slate-50 border border-slate-200 rounded-lg px-3 py-2 text-sm focus:border-blue-400 focus:ring-2 focus:ring-blue-100 outline-none transition-all placeholder-slate-400"
                    placeholder="••••••••"
                  />
                </div>
                <div>
                  <label className="block text-xs text-slate-500 mb-1.5 font-medium">Limit</label>
                  <input
                    type="number"
                    value={limit}
                    onChange={(e) => setLimit(e.target.value)}
                    className="w-full bg-slate-50 border border-slate-200 rounded-lg px-3 py-2 text-sm focus:border-blue-400 focus:ring-2 focus:ring-blue-100 outline-none transition-all"
                    min="1"
                    max="50"
                  />
                </div>
              </div>
            )}
          </div>

          {/* Search Input */}
          <div className="md:col-span-6">
            <div className="relative">
              <span className="absolute left-4 top-1/2 -translate-y-1/2 text-slate-400">
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                  <circle cx="11" cy="11" r="8"></circle>
                  <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
                </svg>
              </span>
              <input
                type="text"
                value={keyword}
                onChange={(e) => setKeyword(e.target.value)}
                className="w-full bg-white border border-slate-200 rounded-xl py-3 pl-11 pr-4 text-sm focus:border-blue-400 focus:ring-2 focus:ring-blue-100 outline-none transition-all placeholder-slate-400 shadow-sm"
                placeholder="Search keyword (e.g. 'Software Engineer', 'Product Manager')"
              />
            </div>

            {/* Message Input */}
            <div className="mt-3">
              <textarea
                value={message}
                onChange={(e) => setMessage(e.target.value)}
                className="w-full bg-white border border-slate-200 rounded-xl py-3 px-4 text-sm focus:border-blue-400 focus:ring-2 focus:ring-blue-100 outline-none transition-all placeholder-slate-400 resize-none shadow-sm"
                rows="2"
                placeholder="Connection message (optional)"
              />
            </div>
          </div>

          {/* Start Button */}
          <div className="md:col-span-3 flex items-start">
            <button
              onClick={handleStart}
              disabled={isRunning || !keyword}
              className={`w-full flex items-center justify-center gap-2 px-6 py-3 rounded-xl font-medium text-sm transition-all ${isRunning
                  ? 'bg-slate-100 text-slate-400 cursor-not-allowed border border-slate-200'
                  : keyword
                    ? 'bg-blue-600 text-white hover:bg-blue-700 shadow-md shadow-blue-200'
                    : 'bg-slate-100 text-slate-400 cursor-not-allowed border border-slate-200'
                }`}
            >
              {isRunning ? (
                <>
                  <svg className="animate-spin w-4 h-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  Running...
                </>
              ) : (
                <>
                  <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="currentColor">
                    <polygon points="5 3 19 12 5 21 5 3"></polygon>
                  </svg>
                  Start
                </>
              )}
            </button>
          </div>
        </div>

        {/* Terminal */}
        <div className="bg-white rounded-2xl border border-slate-200 overflow-hidden shadow-sm">
          {/* Terminal Header */}
          <div className="flex items-center justify-between px-4 py-3 bg-slate-50 border-b border-slate-200">
            <div className="flex items-center gap-2">
              <div className="w-3 h-3 rounded-full bg-red-400"></div>
              <div className="w-3 h-3 rounded-full bg-amber-400"></div>
              <div className="w-3 h-3 rounded-full bg-emerald-400"></div>
            </div>
            <span className="text-xs text-slate-500 font-mono">automation.log</span>
            <div className="flex items-center gap-2 text-xs text-slate-400">
              <span>{logs.length} entries</span>
            </div>
          </div>

          {/* Terminal Body */}
          <div className="h-80 overflow-y-auto p-4 font-mono text-sm custom-scrollbar bg-slate-900">
            {logs.length === 0 && !isRunning && (
              <div className="flex items-center gap-3 text-slate-500">
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className="text-slate-600">
                  <polyline points="4 17 10 11 4 5"></polyline>
                  <line x1="12" y1="19" x2="20" y2="19"></line>
                </svg>
                <span className="text-slate-500">Ready. Enter a keyword and click Start to begin.</span>
              </div>
            )}

            <div className="space-y-1">
              {logs.map((log, i) => (
                <div key={i} className="flex gap-3 py-0.5 hover:bg-slate-800/50 rounded px-2 -mx-2 transition-colors">
                  <span className="text-slate-500 shrink-0 tabular-nums">{log.time}</span>
                  <span className={log.isError ? 'text-red-400' : log.text.includes('successful') ? 'text-emerald-400' : log.text.includes('Skipping') ? 'text-amber-400' : 'text-slate-300'}>{log.text}</span>
                </div>
              ))}
            </div>

            {isRunning && (
              <div className="flex items-center gap-2 mt-2 text-slate-400">
                <span className="animate-pulse">▋</span>
              </div>
            )}
            <div ref={logsEndRef} />
          </div>
        </div>

        {/* Footer */}
        <div className="flex items-center justify-between text-xs text-slate-400 pt-2">
          <span>8 stealth techniques active</span>
          <span className="font-mono">~/.linkedin-automation-profile</span>
        </div>
      </div>
    </div>
  )
}

export default App