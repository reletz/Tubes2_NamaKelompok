import React from 'react';
import tavernSign from '../media/tavern-sign.png';
import logoAndHeads from '../media/logo-and-heads.png';
import log from '../media/Log.png';
import logLight from '../media/Log-light.png';
import magnifier from '../media/magnifier.png';
import '../App.css';
import { Link } from 'react-router-dom';

const Home = () => {
  return (
    <div className="Home-container">
      <div className="Home-title">
        <p>Little&nbsp;&nbsp;Alchemy</p>
        <p>Cookbook</p>
      </div>
      <div className="Home-authors">
        <p>BY LANA, LUCAS, & NAU</p>
      </div>
      <div className="Home-image-container">
        <img src={tavernSign} className="Tavern-sign" alt="Tavern Sign"/>
        <img src={logoAndHeads} className="Logo-and-heads" alt="Logo & Heads"/>
        <Link to="/search">
          <img src={log} className="Log" alt="Log" onMouseOver={e => e.currentTarget.src = logLight} onMouseOut={e => e.currentTarget.src = log} />
        </Link>
        <img src={magnifier} className="Magnifier" alt="Magnifier"/>
      </div>
      <div className="Home-button-search">
        <p>Cari Resep</p>
      </div>
      <div className="Home-button-others">
        <p><Link to="/recipes" className="Home-link-style">Lihat<br/>Resep</Link></p>
        <p><Link to="/ingredients" className="Home-link-style">Lihat<br/>Bahan</Link></p>
        <p><Link to="/about" className="Home-link-style">Tentang<br/>Kami</Link></p>
      </div>
    </div>
  );
};

export default Home;
