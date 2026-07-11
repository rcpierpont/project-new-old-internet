document.addEventListener('DOMContentLoaded', async () => {
  const token = localStorage.getItem('token');

  if (token) {
    document.getElementById('auth-section').style.display = 'none';
    document.getElementById('kreeyaw-create-section').style.display = 'block';
    //document.getElementById('video-section').style.display = 'block';
    //await getVideos();
  } else {
    document.getElementById('auth-section').style.display = 'block';
    document.getElementById('kreeyaw-create-section').style.display = 'none';
    //document.getElementById('video-section').style.display = 'none';
  }
  await logout();
});

document.getElementById('login-form').addEventListener('submit', async (event) => {
  event.preventDefault();
  await login();
  document.getElementById('kreeyaw-create-section').style.display = 'block';
});

document.getElementById('kreeyaw-form').addEventListener('postKreeyaw', async (event) => {
  event.preventDefault();
  console.log("posting kreeyaw");
  await postKreeyaw();
});

async function postKreeyaw() {
  const body = document.getElementById('body').value;
  console.log(`body of kreeyaw is: ${body}`);

  try {
    const res = await fetch('/api/kreeyaws', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
      },
      body: JSON.stringify({ body }),
    });
    
    const data = await res.json();
    if (!res.ok) {
      throw new Error(`Failed to post kreeyaw: ${data.error}`);
    }

    if (data.id) {
      alert('kreeyaw posted successfully!');
    } else {
      throw new Error(`Failed to post kreeyaw: ${data.error}`);
    }

  } catch (error) {
    alert(`Error: ${error.message}`);
  }
}

async function login() {
  const email = document.getElementById('email').value;
  const password = document.getElementById('password').value;

  try {
    const res = await fetch('/api/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email, password }),
    });
    const data = await res.json();
    if (!res.ok) {
      throw new Error(`Failed to login: ${data.error}`);
    }

    if (data.token) {
      localStorage.setItem('token', data.token);
      document.getElementById('auth-section').style.display = 'none';
    } else {
      alert('Login failed. Please check your credentials.');
    }
  } catch (error) {
    alert(`Error: ${error.message}`);
  }
}

async function signup() {
  const email = document.getElementById('email').value;
  const password = document.getElementById('password').value;

  try {
    const res = await fetch('/api/users', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email, password }),
    });
    if (!res.ok) {
      const data = await res.json();
      throw new Error(`Failed to create user: ${data.error}`);
    }
    console.log('User created!');
    await login();
  } catch (error) {
    alert(`Error: ${error.message}`);
  }
}

function logout() {
  localStorage.removeItem('token');
  document.getElementById('auth-section').style.display = 'block';
}
