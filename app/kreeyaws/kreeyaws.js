document.addEventListener('DOMContentLoaded', async () => {
  const token = localStorage.getItem('token');

  if (token) {
    document.getElementById('kreeyaw-section').style.display = 'none';
    //document.getElementById('video-section').style.display = 'block';
    //await getVideos();
  } else {
    document.getElementById('nothing-section').style.display = 'none';
    //document.getElementById('video-section').style.display = 'none';
  }
  await logout();
});