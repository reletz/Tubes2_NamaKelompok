import React from 'react';
import { BrowserRouter, Routes, Route, Link } from 'react-router-dom';
import Home from './pages/Home';
import About from './pages/About';
import Search from './pages/Search';
import Recipes from './pages/Recipes';
import Ingredients from './pages/Ingredients';

const App = () => {
  return (
    <BrowserRouter>
      <nav>
        <Link to="/">Home</Link> |
        <Link to="/about">About</Link> |
        <Link to="/search">Search</Link> |
        <Link to="/recipes">Recipes</Link> |
        <Link to="/ingredients">Ingredients</Link> |
      </nav>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/about" element={<About />} />
        <Route path="/search" element={<Search />} />
        <Route path="/recipes" element={<Recipes />} />
        <Route path="/ingredients" element={<Ingredients />} />
      </Routes>
    </BrowserRouter>
  );
};

export default App;
