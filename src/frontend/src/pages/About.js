import React from 'react';
import nau from '../media/icons/nau.png';
import lucas from '../media/icons/lucas.png';
import lana from '../media/icons/lana.png';
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
            <p>Samantha Laqueenna Ginting, akrab disapa Lana, memegang peranan penting dalam tugas besar ini. Sebagai Frontend Developer, ia bertanggung jawab dalam merancang dan mengimplementasikan antarmuka yang responsif dan interaktif. Lana menggunakan React.js sebagai kerangka kerja utama, memastikan bahwa setiap komponen visual bekerja secara mulus dan terintegrasi dengan baik dengan backend. Ia juga berperan dalam menyusun alur interaksi pengguna, mengatur state management menggunakan React Hooks, serta menghubungkan API dari backend agar data dapat ditampilkan secara real-time.</p>
          </div>
        </div>

        <div className="About-card">
          <img src={lucas} className="About-avatar" alt="Lucas Avatar" />
          <div className="About-info">
            <h2>NICHOLAS ANDHIKA LUCAS</h2>
            <h3>NIM 13523014</h3>
            <p>Nicholas Andhika Lucas, akrab disapa Lucas, memegang peranan penting dalam tugas besar ni. Sebagai UI/UX Designer dan Frontend Developer, Lucas bertanggung jawab dalam merancang tampilan antarmuka yang tidak hanya menarik, tetapi juga mudah digunakan dan intuitif bagi pengguna. Dalam perannya sebagai Frontend Developer, Lucas mengimplementasikan desain tersebut ke dalam kode menggunakan framework seperti React.js, memastikan setiap elemen UI berfungsi dengan baik di berbagai perangkat dan ukuran layar.</p>
          </div>
        </div>

        <div className="About-card">
          <img src={nau} className="About-avatar" alt="Nau Avatar" />
          <div className="About-info">
            <h2>NAUFARREL ZHAFIF ABHISTA</h2>
            <h3>NIM 13523149</h3>
            <p>Naufarrel Zhafif Abhista, akrab disapa Nau, memegang peranan penting dalam tugas besar ni. Sebagai Backend Developer, Nau bertanggung jawab membangun dan mengelola logika serta struktur data di sisi server. Dalam tugas besar ini, Nau juga berperan penting dalam merancang dan mengimplementasikan algoritma pencarian seperti BFS dan DFS yang digunakan untuk menemukan jalur pembuatan resep dalam sistem. Ia mengoptimalkan algoritma tersebut agar dapat bekerja secara cepat meskipun pada dataset yang besar dan kompleks.</p>
          </div>
        </div>

      </div>

    </div>
  );
};

export default About;
