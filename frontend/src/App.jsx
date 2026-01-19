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
        setLogs((prev) => [...prev, event.data])
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
        throw new Error('Failed to start')
      }
    } catch (err) {
      setLogs(prev => [...prev, `Error: ${err.message}`])
      setIsRunning(false)
    }
  }

  return (
    <div className="min-h-screen bg-[#0a0a0a] text-gray-300 font-mono p-8 flex flex-col items-center">
      {/* Top Control Bar */}
      <div className="w-full max-w-5xl flex flex-wrap items-center justify-between gap-4 mb-8">
        
        {/* Credentials Toggle */}
        <div className="relative">
          <button 
            onClick={() => setShowCreds(!showCreds)}
            className="flex items-center gap-2 px-4 py-2 bg-[#111] border border-gray-800 rounded-md hover:border-gray-600 transition-colors text-sm"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect><path d="M7 11V7a5 5 0 0 1 10 0v4"></path></svg>
            {showCreds ? 'Hide Credentials' : 'Custom Credentials (Optional)'}
          </button>

          {showCreds && (
            <div className="absolute top-12 left-0 z-10 w-72 bg-[#111] border border-gray-800 rounded-md p-4 shadow-xl">
              <div className="space-y-3">
                <div>
                  <label className="block text-xs text-gray-500 mb-1">Email</label>
                  <input 
                    type="email" 
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    className="w-full bg-black border border-gray-800 rounded px-2 py-1 text-sm focus:border-gray-500 outline-none"
                    placeholder="user@example.com"
                  />
                </div>
                <div>
                  <label className="block text-xs text-gray-500 mb-1">Password</label>
                  <input 
                    type="password" 
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    className="w-full bg-black border border-gray-800 rounded px-2 py-1 text-sm focus:border-gray-500 outline-none"
                    placeholder="••••••••"
                  />
                </div>
                <div>
                  <label className="block text-xs text-gray-500 mb-1">Limit</label>
                  <input 
                    type="number" 
                    value={limit}
                    onChange={(e) => setLimit(e.target.value)}
                    className="w-full bg-black border border-gray-800 rounded px-2 py-1 text-sm focus:border-gray-500 outline-none"
                    min="1"
                  />
                </div>
              </div>
            </div>
          )}
        </div>

        {/* Search Bar */}
        <div className="flex-1 max-w-xl">
          <div className="relative">
            <span className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>
            </span>
            <input 
              type="text" 
              value={keyword}
              onChange={(e) => setKeyword(e.target.value)}
              className="w-full bg-[#111] border-none rounded-md py-2.5 pl-10 pr-4 text-sm focus:ring-1 focus:ring-gray-700 outline-none placeholder-gray-600"
              placeholder="Search Bar (e.g. 'Software Engineer')"
            />
          </div>
        </div>

        {/* Start Button */}
        <button 
          onClick={handleStart}
          disabled={isRunning || !keyword}
          className={`flex items-center gap-2 px-6 py-2.5 rounded-md font-medium text-sm transition-all ${
            isRunning 
              ? 'bg-gray-800 text-gray-500 cursor-not-allowed' 
              : 'bg-white text-black hover:bg-gray-200'
          }`}
        >
          {isRunning ? (
            <>
              <span className="animate-spin">⟳</span> Running...
            </>
          ) : (
            <>
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor" stroke="none"><polygon points="5 3 19 12 5 21 5 3"></polygon></svg>
              Start Automation
            </>
          )}
        </button>
      </div>

      {/* Terminal Output */}
      <div className="w-full max-w-5xl flex-1 bg-[#0f0f0f] rounded-lg border border-gray-800/50 p-6 overflow-hidden flex flex-col shadow-2xl">
        <div className="flex gap-2 mb-4">
          <div className="w-3 h-3 rounded-full bg-gray-800"></div>
          <div className="w-3 h-3 rounded-full bg-gray-800"></div>
          <div className="w-3 h-3 rounded-full bg-gray-800"></div>
        </div>
        
        <div className="flex-1 overflow-y-auto font-mono text-sm space-y-2 custom-scrollbar">
          {logs.length === 0 && !isRunning && (
            <div className="text-gray-600 italic">
              Ready to start. Enter a keyword and click start.
            </div>
          )}
          
          {logs.map((log, i) => (
            <div key={i} className="text-gray-300">
              <span className="text-gray-600 mr-3">{new Date().toLocaleTimeString()}</span>
              {log}
            </div>
          ))}
          
          {isRunning && (
            <div className="animate-pulse text-gray-500">_</div>
          )}
          <div ref={logsEndRef} />
        </div>
      </div>
    </div>
  )
}

export default App