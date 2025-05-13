import React from 'react';
import nau from '../media/nau.png';
import lucas from '../media/lucas.png';
import lana from '../media/lana.png';
import '../App.css';

const About = () => {
  return (
    <div className="About-container">
      <h1 className="About-title">TENTANG&nbsp;&nbsp;KAMI</h1>

      <div className="About-description">
        <p>
        Little Alchemy Cookbook adalah sebuah aplikasi berbasis website untuk melakukan pencarian resep elemen dalam permainan Little Alchemy 2 dengan menggunakan strategi BFS dan DFS. Frontend website ini dibangun menggunakan framework React.js, sementara Backend dibangun mengunakan bahasa Go. Data seluruh elemen dan resep diperoleh dari website <a href='https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)'>little-alchemy.fandom.com</a>. Mode pencarian multiple recipe dioptimasi dengan multithreading.
        </p>
        <h2>TIM&nbsp;&nbsp;PENGEMBANG</h2>
      </div>

      <div className="About-cards">

        <div className="About-card">
          <img src={lana} className="About-avatar" alt="Lana Avatar" />
          <div className="About-info">
            <h2>SAMANTHA LAQUEENNA GINTING</h2>
            <h3>NIM 13523138</h3>
            <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nam ultrices turpis eu libero elementum, ut cursus turpis commodo. Vestibulum dui arcu, aliquet vitae mauris vel, congue mollis nunc. Etiam quam mi, convallis ac elit id, sagittis cursus eros.</p>
          </div>
        </div>

        <div className="About-card">
          <img src={lucas} className="About-avatar" alt="Lucas Avatar" />
          <div className="About-info">
            <h2>NICHOLAS ANDHIKA LUCAS</h2>
            <h3>NIM 13523014</h3>
            <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nam ultrices turpis eu libero elementum, ut cursus turpis commodo. Vestibulum dui arcu, aliquet vitae mauris vel, congue mollis nunc. Etiam quam mi, convallis ac elit id, sagittis cursus eros.</p>
          </div>
        </div>

        <div className="About-card">
          <img src={nau} className="About-avatar" alt="Nau Avatar" />
          <div className="About-info">
            <h2>NAUFARREL ZHAFIF ABHISTA</h2>
            <h3>NIM 13523149</h3>
            <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nam ultrices turpis eu libero elementum, ut cursus turpis commodo. Vestibulum dui arcu, aliquet vitae mauris vel, congue mollis nunc. Etiam quam mi, convallis ac elit id, sagittis cursus eros.</p>
          </div>
        </div>

      </div>

    </div>
  );
};

export default About;
