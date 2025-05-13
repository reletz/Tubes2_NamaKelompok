import { React, useState, useRef, useEffect } from 'react';
import { useForm } from 'react-hook-form';
import Tree from 'react-d3-tree';
import telenanBg from '../media/talenan.png';
import BFS from '../media/icons/BFS.svg';
import DFS from '../media/icons/DFS.svg';
import single from '../media/icons/single.svg';
import multiple from '../media/icons/multiple.svg';
import biBFS from '../media/icons/biBFS.svg';

const About = () => {
  const { register, handleSubmit } = useForm();
  const [treeDataList, setTreeDataList] = useState([]);
  const [metaInfo, setMetaInfo] = useState({ timetaken: "-", node_visited: 0 });
  const treeContainerRef = useRef(null);
  const [containerWidth, setContainerWidth] = useState(0);
  const [selectedAlgorithm, setSelectedAlgorithm] = useState(null);
  const [selectedSearchMode, setSelectedSearchMode] = useState(null);

  const onSubmit = async (querySearch) => {
    querySearch.maksimalResep = Number(querySearch.maksimalResep);
    console.log("Kirim ke backend:", querySearch); // debug di console (liat inspect)

    try {
      const response = await fetch("http://localhost:8080/api/search", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(querySearch),
      });

      const data = await response.json();

      if (data.treeData && Array.isArray(data.treeData)) {
        setTreeDataList(data.treeData);
        setMetaInfo({
          timetaken: data.timetaken || "-",
          node_visited: data.node_visited || 0,
        });
      } else {
        console.error("Unexpected API response format:", data);
      }
    } catch (err) {
      console.error("Failed to fetch:", err);
    }
  };

  useEffect(() => {
    if (treeContainerRef.current) {
      setContainerWidth(treeContainerRef.current.offsetWidth);
    }
  }, [treeDataList]);

  const calculateTreeDepth = (node) => {
    if (!node.children || node.children.length === 0) return 1;
    return 1 + Math.max(...node.children.map(calculateTreeDepth));
  };

  const calculateTreeWidth = (node, levelWidths = {}, depth = 0) => {
    if (!levelWidths[depth]) levelWidths[depth] = 0;
    levelWidths[depth] += 1;
    if (node.children) {
      node.children.forEach(child => calculateTreeWidth(child, levelWidths, depth + 1));
    }
    return Math.max(...Object.values(levelWidths));
  };
  const renderCustomNode = ({ nodeDatum }) => {
    const isLeaf = !nodeDatum.children || nodeDatum.children.length === 0;
    const textLength = nodeDatum.name.length;
    const width = 70;
    const baseHeight = 35;
    const estimatedLineCount = Math.ceil(textLength / 16);
    const height = baseHeight * estimatedLineCount;
    return (
      <g>
        <foreignObject width={width} height={height} x={-width / 2} y={-height / 2}>
          <div
            xmlns="http://www.w3.org/1999/xhtml"
            className={`tree-node ${isLeaf ? 'leaf-node' : 'parent-node'}`}
          >
            <p className="node-label-multiline">{nodeDatum.name}</p>
          </div>
        </foreignObject>
      </g>
    );
  };

  const OptionsButton = ({ options, dataType, selectedOption, setSelectedOption }) => {
    const handleClick = (value) => {
      setSelectedOption(value);
    };

    return (
    <div className="options-container">
      {options.map((option, index) => (
        <button
          key={index}
          type="button"
          className={`options-button ${selectedOption === option.name ? 'active' : ''}`}
          onClick={() => handleClick(option.name)}
        >
          <img src={option.icon} alt={option.name} />
          {option.name}
        </button>
      ))}

      <input
        type="hidden"
        value={selectedOption || ""}
        {...register(dataType, { required: true })}
      />
    </div>
  );
};

  return (
    <div className="Search-container">
      <h1 className="Search-title">CARI&nbsp;&nbsp;&nbsp;RESEP</h1>

      <form onSubmit={handleSubmit(onSubmit)}>
        <div className="Search-form">
          <div className="Search-form-card">
            <h3>Nama Resep*</h3>
            <input
              type="text"
              placeholder="Contoh: Brick, Cloud"
              autoComplete="off"
              className="custom-search-input"
              {...register("namaResep", { required: true })}
            />
          </div>

          <div className="Search-form-card">
            <h3>Jumlah Resep*</h3>
            <input
              type="number"
              placeholder="Contoh: 5"
              className="custom-search-input"
              {...register("maksimalResep", { required: true, min: 1 })}
            />
          </div>

          <div className="Search-form-card">
            <h3>Algoritma*</h3>
            <OptionsButton
              options={[
                { name: 'BFS', icon: BFS },
                { name: 'DFS', icon: DFS },
                { name: 'Bi-BFS', icon: biBFS },
              ]}
              selectedOption={selectedAlgorithm}
              setSelectedOption={setSelectedAlgorithm}
              dataType="algoritma"
            />
          </div>

          <div className="Search-form-card">
            <h3>Mode Pencarian*</h3>
            <OptionsButton
              options={[
                { name: 'Single', icon: single },
                { name: 'Multiple', icon: multiple },
              ]}
              selectedOption = {selectedSearchMode}
              setSelectedOption={setSelectedSearchMode}
              dataType="modePencarian"
            />
          </div>       

          <input type="submit" className="submit-button" />
        </div>
      </form>

      {/* Tree View */}
      {treeDataList.length > 0 && (
        <div className="Search-Tree" ref={treeContainerRef} style={{ width: '100%', minHeight: '100vh' }}>
          <h1 className="Search-title">HASIL PENCARIAN</h1>
          <div className="Search-meta-info">
            <p>Waktu Pencarian: {metaInfo.timetaken}</p>
            <p>Node yang Dikunjungi: {metaInfo.node_visited}</p>
          </div>
          {treeDataList.map((treeData, idx) => {
            const depth = calculateTreeDepth(treeData);
            const height = Math.max(depth * 80, 120);
            const treeWidth = calculateTreeWidth(treeData) * 120;
            return (
              <div key={idx} style={{ height: `${height}px`}}>
                <h3 className="Search-subtitle-2">Resep {treeData.name} #{idx+1}</h3>
                <div
                  style={{
                    width: '80%',
                    margin: '0 auto',
                    height: `${height}px`,
                    backgroundImage: `url(${telenanBg})`,
                    backgroundSize: 'cover',
                    backgroundRepeat: 'no-repeat',
                    backgroundPosition: 'center',
                  }}
                >
                  <Tree
                    data={treeData}
                    orientation="vertical"
                    translate={{ x: containerWidth * 0.4, y: 50 }}
                    nodeSize={{ x: 65, y: 60 }}
                    zoomable={true}
                    initialZoom={Math.min(containerWidth / treeWidth, 0.3)}
                    pathFunc="diagonal"
                    separation={{ siblings: 1.2, nonSiblings: 1.5 }}
                    renderCustomNodeElement={renderCustomNode}
                  />
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
};

export default About;
