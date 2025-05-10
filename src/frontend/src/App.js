import React from 'react';
import { BrowserRouter, Routes, Route, Link, useLocation } from 'react-router-dom';
import Home from './pages/Home';
import About from './pages/About';
import Search from './pages/Search';
import Ingredients from './pages/Ingredients';

const App = () => {
  return (
    <BrowserRouter>
      <div className="App">
        <AppContent />
      </div>
    </BrowserRouter>
  );
};

const AppContent = () => {
  const location = useLocation();

  return (
    <>
      {location.pathname !== '/' && (
        <nav className="Main-navbar">
          <Link to="/">BERANDA</Link>
          <Link to="/search">CARI RESEP</Link>
          <Link to="/ingredients">LIHAT BAHAN</Link>
          <Link to="/about">TENTANG KAMI</Link>
        </nav>
      )}
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/about" element={<About />} />
        <Route path="/search" element={<Search />} />
        <Route path="/ingredients" element={<Ingredients />} />
      </Routes>
    </>
  );
};

export default App;
