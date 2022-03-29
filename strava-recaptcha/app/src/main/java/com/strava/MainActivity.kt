package com.strava

import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle
import android.widget.Button
import android.widget.TextView
import com.google.android.gms.common.api.ApiException
import com.google.android.gms.common.api.CommonStatusCodes
import com.google.android.gms.safetynet.SafetyNet;

class MainActivity : AppCompatActivity() {
    lateinit var tokenText: TextView;

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)

        tokenText = findViewById(R.id.tokenTextView)

        val btn: Button = findViewById(R.id.captchaButton)
        btn.setOnClickListener { RunReCAPTCHA() }
    }

    fun RunReCAPTCHA() {
        SafetyNet.getClient(this).verifyWithRecaptcha(SAFETY_NET_API_KEY)
            .addOnSuccessListener(this) { response ->
                if (!response.tokenResult.isEmpty()) {
                    tokenText.text = response.tokenResult
                }
            }
            .addOnFailureListener(this) { e ->
                if (e is ApiException) {
                    tokenText.text = "SafetyNet Error: " + CommonStatusCodes.getStatusCodeString(e.statusCode)
                } else {
                    tokenText.text = "Unknown SafetyNet error: " + e.message
                }
            }
    }

    companion object {
        private val SAFETY_NET_API_KEY = "6LcyRkUUAAAAANqIvWNo-kOBWY5yWbboMC_SV8n2"
    }

}