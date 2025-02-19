import React from 'react'
// import Navigation from '../components/Navigation'
// import Footer from '../components/Footer'

const Layout = ({ children }) => {
    return (
        <>
            <main className='my-8 container'>{children}</main>
            {/* <Footer /> */}
        </>
    )
}

export default Layout