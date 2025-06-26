import { useState } from 'react'

const UrlForm = () => {
    const [longUrl, setLongUrl] = useState('')
    const [shortUrl, setShortUrl] = useState('')
    const [error, setError] = useState('')
    const [loading, setLoading] = useState(false)
    const [forceNew, setForceNew] = useState(false)

    const isValidUrl = (url) => {
        try {
            new URL(url)
            return true
        } catch {
            return false
        }
    }

    const handleSubmit = async (e) => {
        e.preventDefault()
        setError('')
        setShortUrl('')

        if (!isValidUrl(longUrl)) {
            setError('Please enter a valid URL starting with http:// or https://')
            return
        }

        setLoading(true)
        try {
            const res = await fetch('/links/shorten', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ long_url: longUrl, force_new: forceNew })
            })

            const data = await res.json()
            if (!res.ok) {
                throw new Error(data?.error || 'Server error')
            }

            setShortUrl(data.short_url)
        } catch (err) {
            setError(err.message)
        } finally {
            setLoading(false)
        }
    }

    return (
        <div style={{ width: '100%', maxWidth: '700px' }}>
            <form
                onSubmit={handleSubmit}
                style={{
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    gap: '1rem'
                }}
            >
                <input
                    type="url"
                    placeholder="Enter long URL..."
                    value={longUrl}
                    onChange={(e) => setLongUrl(e.target.value)}
                    required
                    style={{
                        flexGrow: 1,
                        padding: '0.6rem',
                        border: '1px solid #444',
                        borderRadius: '4px',
                        backgroundColor: '#1e1e1e',
                        color: '#fff'
                    }}
                />

                <button
                    type="submit"
                    style={{
                        padding: '0.6rem 1rem',
                        backgroundColor: '#1e88e5',
                        color: 'white',
                        border: 'none',
                        borderRadius: '4px',
                        cursor: 'pointer',
                        whiteSpace: 'nowrap'
                    }}
                >
                    {loading ? 'Shortening...' : 'Shorten'}
                </button>
            </form>

            <label
                style={{
                    color: '#ccc',
                    marginTop: '1rem',
                    display: 'block',
                    fontSize: '0.95rem'
                }}
            >
                <input
                    type="checkbox"
                    checked={forceNew}
                    onChange={() => setForceNew(!forceNew)}
                    style={{ marginRight: '0.5rem' }}
                />
                Force new short URL
            </label>

            {shortUrl && (
                <p style={{ marginTop: '1rem', color: '#00e676' }}>
                    Short URL: <a href={shortUrl} target="_blank" rel="noreferrer" style={{ color: '#80d8ff' }}>{shortUrl}</a>
                </p>
            )}

            {error && (
                <p style={{ marginTop: '1rem', color: '#ff5252' }}>
                    ⚠️ {error}
                </p>
            )}
        </div>
    )
}

export default UrlForm
