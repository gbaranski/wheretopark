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
    let parkingLotID: String
    let image: UIImage
    let email: String = "contact@wheretopark.app"
    @Environment(\.presentationMode) var presentation
    @Binding var result: Result<MFMailComposeResult, Error>?

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
//        let imageString = image.pngData()!.base64EncodedString(options: .lineLength64Characters)
        vc.setMessageBody(
            """
            <p>Hi, I've got a problem with parking lot of ID: \(parkingLotID)</p>
            <br/>
            <p>My problem is: (describe your problem here)</p>
            """,
//            <img src='data:image/png;base64,\(imageString)' width='\(image.size.width)' height='\(image.size.height)'>
            isHTML: true
        )
        vc.addAttachmentData(image.pngData()!, mimeType: "image/png", fileName: "screenshot.png")
        vc.mailComposeDelegate = context.coordinator
        return vc
    }

    func updateUIViewController(_ uiViewController: MFMailComposeViewController,
                                context: UIViewControllerRepresentableContext<MailView>) {

    }
}
