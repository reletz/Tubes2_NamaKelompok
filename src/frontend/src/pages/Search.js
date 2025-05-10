import React from 'react';
import { useForm } from 'react-hook-form';

const About = () => {
  const { register, handleSubmit, formState: { errors } } = useForm();
  const onSubmit = querySearch => console.log(querySearch); // Query adalah variabel yang bentuknya menyerupai JSON
  console.log(errors);

  return (
          <div className="Search-container">
            <h1 className="Search-title">CARI&nbsp;&nbsp;&nbsp;RESEP</h1>
            
              <form onSubmit={handleSubmit(onSubmit)}>

                <div className="Search-form">
                  <div className="Search-form-card">
                    <h3>Nama Resep*</h3>
                    <input
                      type="text"
                      placeholder="Contoh: Babe the blue ox"
                      autoComplete='off'
                      className="custom-search-input"
                      {...register("Nama Resep", { required: true })}
                    />
                  </div>

                  <div className="Search-form-card">
                    <h3>Jumlah Resep*</h3>
                    <input 
                      type="number"
                      placeholder="Contoh: 5"
                      className="custom-search-input"
                      {...register("Maksimal Resep", {required: true, min: 1})} />
                  </div>

                  <div className="Search-form-card">
                    <h3>Algoritma*</h3>
                    <div className="Search-form-card-radio">
                    <div className="Search-form-card-radio-2">
                        <label className="custom-radio">
                          <input
                            type="radio"
                            {...register("Algoritma", { required: true })}
                            value="BFS"
                          />
                          <span className="radio-image" />
                        </label>
                        
                        <label className="custom-radio">
                          <input
                            type="radio"
                            {...register("Algoritma", { required: true })}
                            value="DFS"
                          />
                          <span className="radio-image" />
                        </label>
                      </div>
                      <div className="Search-form-card-radio-3">
                        <p>BFS</p>
                        <p>DFS</p>
                      </div>
                    </div>
                  </div>

                  <input 
                  type="submit" 
                  className="submit-button" />

                </div>
              </form>
          </div>
        );
};

export default About;