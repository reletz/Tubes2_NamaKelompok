import React from 'react';
import { useForm } from 'react-hook-form';

const About = () => {
  const { register, handleSubmit, formState: { errors } } = useForm();
  const onSubmit = query => console.log(query); // Query adalah variabel yang bentuknya menyerupai JSON
  console.log(errors);

  return (
          <div className="Search-container">
            <h1 className="Search-title">CARI&nbsp;&nbsp;RESEP</h1>
            
              <form onSubmit={handleSubmit(onSubmit)}>

                <div className="Search-form">
                  <div className="Search-form-card">
                    <h3>Nama Resep</h3>
                    <input type="text" placeholder=" " {...register("Nama Resep", {required: true})} />
                  </div>

                  <div className="Search-form-card">
                    <h3>Algoritma</h3>
                    <div className="Search-form-card-radio">
                      <div className="Search-form-card-radio-2">
                        <input {...register("Algoritma", { required: true })} type="radio" value="BFS" />
                        <input {...register("Algoritma", { required: true })} type="radio" value="DFS" />
                      </div>
                      <div className="Search-form-card-radio-3">
                        <p>BFS</p>
                        <p>DFS</p>
                      </div>
                    </div>
                  </div>

                  <div className="Search-form-card">
                    <h3>Mode Pencarian</h3>
                    <div className="Search-form-card-radio">
                      <div className="Search-form-card-radio-2">
                        <input {...register("Mode", { required: true })} type="radio" value="Resep Terpendek" />
                        <input {...register("Mode", { required: true })} type="radio" value="Lihat Banyak Resep" />
                      </div>
                      <div className="Search-form-card-radio-3">
                        <p>Resep Terpendek</p>
                        <p>Lihat Banyak Resep</p>
                      </div>
                    </div>
                  </div>

                  <div className="Search-form-card">
                    <h3>Jumlah Resep</h3>
                    <input type="number" placeholder=" " {...register("Maksimal Resep", {required: true, min: 1})} />
                  </div>

                  <input type="submit" />

                </div>

              </form>
            
          </div>
        );
};

export default About;
