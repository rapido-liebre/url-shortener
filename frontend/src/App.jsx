import UrlForm from './components/UrlForm'

function App() {
    return (
        <div
            style={{
                minHeight: '100vh',
                backgroundColor: '#121212',
                color: '#fff',
                paddingTop: '3rem',
                display: 'flex',
                flexDirection: 'column',
                alignItems: 'center',
                fontFamily: 'Segoe UI, sans-serif'
            }}
        >
            <h1 style={{ marginBottom: '2rem' }}>🔗 URL Shortener</h1>
            <UrlForm />
        </div>
    )
}

export default App
