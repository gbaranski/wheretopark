//
//  MailView.swift
//  wheretopark
//
//  Created by Grzegorz Barański on 13/09/2022.
//

import Foundation
import SwiftUI
import UIKit
import MessageUI

struct MailView: UIViewControllerRepresentable {
    let message: () -> String
    let attachment: (() -> UIImage)?
    @Binding var result: Result<MFMailComposeResult, Error>?
    
    let email: String = "contact@wheretopark.app"
    @Environment(\.presentationMode) var presentation

    class Coordinator: NSObject, MFMailComposeViewControllerDelegate {

        @Binding var presentation: PresentationMode
        @Binding var result: Result<MFMailComposeResult, Error>?

        init(presentation: Binding<PresentationMode>,
             result: Binding<Result<MFMailComposeResult, Error>?>) {
            _presentation = presentation
            _result = result
        }

        func mailComposeController(_ controller: MFMailComposeViewController,
                                   didFinishWith result: MFMailComposeResult,
                                   error: Error?) {
            defer {
                $presentation.wrappedValue.dismiss()
            }
            guard error == nil else {
                self.result = .failure(error!)
                return
            }
            self.result = .success(result)
        }
    }

    func makeCoordinator() -> Coordinator {
        return Coordinator(presentation: presentation,
                           result: $result)
    }

    func makeUIViewController(context: UIViewControllerRepresentableContext<MailView>) -> MFMailComposeViewController {
        let vc = MFMailComposeViewController()
        vc.setToRecipients([email])
        vc.setSubject("User feedback")
        vc.setMessageBody(message(), isHTML: true)
        if let attachment = attachment?() {
            vc.addAttachmentData(attachment.pngData()!, mimeType: "image/png", fileName: "screenshot.png")
        }
        vc.mailComposeDelegate = context.coordinator
        return vc
    }

    func updateUIViewController(_ uiViewController: MFMailComposeViewController,
                                context: UIViewControllerRepresentableContext<MailView>) {

    }
}


struct SendFeedback: View {
    let message: () -> String
    let attachment: (() -> UIImage)?
    
    @State private var result: Result<MFMailComposeResult, Error>? = nil
    @State private var isShowingMailView = false

    var body: some View {
        Button(action: {
            self.isShowingMailView = true
        }) {
            Text("Send a feedback")
        }
        .disabled(!MFMailComposeViewController.canSendMail())
        .sheet(isPresented: $isShowingMailView) {
            MailView(
                message: message,
                attachment: attachment,
                result: self.$result
            )
        }
    }
}
