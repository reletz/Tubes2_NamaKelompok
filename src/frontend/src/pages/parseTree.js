export function parseMultipleTrees(rawData) {
  if (!rawData || !Array.isArray(rawData.recipes)) return [];
  return rawData.recipes.map((tree, index) => ({
    ...tree,
    name: tree.name + " #" + (index + 1)
  }));
}
export function parseMetaInfo(rawData) {
  if (!rawData) return { timetaken: "-", node_visited: 0 };
  return {
    timetaken: rawData.timetaken || "-",
    node_visited: rawData.node_visited || 0
  };
}