import React from 'react';
import alembic from '../alembic.svg';
import '../App.css';

const Home = () => {
    return (
        <div className="App">
          <header className="App-header">
            <img src={alembic} className="App-logo" alt="logo" />
            <p>
              NamaKelompok Recipe Finder
            </p>
            <a
              className="App-link"
              href='./about'
            //   target="_blank"
            //   rel="noopener noreferrer"
            >
              Search Recipe
            </a>
          </header>
        </div>
      );
};

export default Home;
