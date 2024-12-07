async function custom (c) {
    const {msg= "", title= ""} = c;
    const {value: result} = Swal.fire({
        title: title,
        html: msg,
        backdrop: false,
        focusConfirm: false,
        showCancelButton: true,
        showConfirmButton: true,  // Делаем кнопку подтверждения
        confirmButtonText: 'Confirm', // Текст подтверждения
        confirmButtonColor: '#3085d6',
        willOpen: () => {
            document.getElementById("submit-button").addEventListener("click", () => {
                Swal.clickConfirm(); // Программно "нажимаем" кнопку подтверждения
            });

            const elem =  document.getElementById('reservation-dates-modal');
            const rp = new Date().getTime();
        },
        preConfirm: async (login) => {
            return [
                document.getElementById('start').value,
                document.getElementById('end').value,
            ]
        },
        didOpen: () => {
            document.getElementById('start').removeAttribute('disabled');
            document.getElementById('end').removeAttribute('disabled');
        }
    });

    if (result) {
        console.log(result)
        if (result.dismiss !== Swal.DismissReason.cancel) {
            if (result.value !== "") {
                if (c.callback !== undefined) {
                    c.callback(result);
                }
            } else {
                c.callback(false)
            }
        } else {
            c.callback(false)
        }
    }
}