document.addEventListener('DOMContentLoaded', async () => {
  const token = localStorage.getItem('token');
  
  if (token) {
    const kreeyaws = await getKreeyaws();
    if (kreeyaws.length > 0) {
      document.getElementById('kreeyaw-section').style.display = 'block';
    } 
  } else {
    window.location.replace("http://localhost:8080/");
  }
  await logout();
});

async function logout() {
  localStorage.removeItem('token');
  window.location.replace("http://localhost:8080/");
}