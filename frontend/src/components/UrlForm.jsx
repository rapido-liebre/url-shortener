import { useState } from 'react'

const UrlForm = () => {
    const [longUrl, setLongUrl] = useState('')
    const [shortUrl, setShortUrl] = useState('')
    const [error, setError] = useState('')

    const handleSubmit = async (e) => {
        e.preventDefault()
        setShortUrl('')
        setError('')

        try {
            const res = await fetch('/links/shorten', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ long_url: longUrl })
            })

            if (!res.ok) throw new Error('Invalid URL or server error')
            const data = await res.json()
            setShortUrl(data.short_url)
        } catch (err) {
            setError(err.message)
        }
    }

    return (
        <>
            <form onSubmit={handleSubmit}>
                <input
                    type="url"
                    placeholder="Enter long URL..."
                    value={longUrl}
                    onChange={(e) => setLongUrl(e.target.value)}
                    required
                    style={{ width: '300px', marginRight: '1rem' }}
                />
                <button type="submit">Shorten</button>
            </form>
            {shortUrl && (
                <p>
                    Short URL: <a href={shortUrl} target="_blank" rel="noreferrer">{shortUrl}</a>
                </p>
            )}
            {error && <p style={{ color: 'red' }}>{error}</p>}
        </>
    )
}

export default UrlForm
