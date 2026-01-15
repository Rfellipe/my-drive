import { NavLink, Route, Routes } from 'react-router'
import './App.css'
import Home from './pages/Home'
import About from './pages/About'

function App() {
  return (
    <>
      <header>
        <nav>
          <NavLink to="/">home</NavLink>
          <NavLink to="/about">about</NavLink>
        </nav>
      </header>

      <main>
        <Routes>
          <Route path='/' element={<Home />} />
          <Route path='/about' element={<About />} />
        </Routes>
      </main>
    </>

  )
}

export default App
